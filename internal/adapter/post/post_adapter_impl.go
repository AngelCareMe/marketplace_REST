package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"marketplace/internal/entity"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostAdapter struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewPostAdapter(db *pgxpool.Pool, logger *logrus.Logger) *PostAdapter {
	return &PostAdapter{
		db:     db,
		logger: logger,
	}
}

func (a *PostAdapter) Create(ctx context.Context, post *entity.Post) error {
	// Проверка на уникальность поста
	queryCheck, argsCheck, err := squirrel.Select("id").
		From("posts").
		Where(squirrel.Eq{
			"header":    post.Header,
			"content":   post.Content,
			"author_id": post.AuthorID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build check post query")
		return fmt.Errorf("check post existence: %w", err)
	}
	var existingID uuid.UUID
	err = a.db.QueryRow(ctx, queryCheck, argsCheck...).Scan(&existingID)
	if err == nil {
		a.logger.WithFields(logrus.Fields{
			"header":    post.Header,
			"author_id": post.AuthorID,
		}).Warn("Post with the same header, content, and author already exists")
		return fmt.Errorf("post with the same header, content, and author already exists")
	}
	if err != sql.ErrNoRows {
		a.logger.WithError(err).Error("Failed to check post existence")
		return fmt.Errorf("check post existence: %w", err)
	}

	// Создание поста
	query, args, err := squirrel.Insert("posts").
		Columns("id", "header", "content", "image", "price", "author_id", "created_at").
		Values(post.ID, post.Header, post.Content, post.Image, post.Price, post.AuthorID, post.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build create post query")
		return fmt.Errorf("create post query: %w", err)
	}
	_, err = a.db.Exec(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to create post")
		return fmt.Errorf("create post: %w", err)
	}
	a.logger.WithFields(logrus.Fields{
		"post_id":   post.ID,
		"author_id": post.AuthorID,
	}).Info("Post created in database")
	return nil
}

func (a *PostAdapter) GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	query, args, err := squirrel.Select("p.id", "p.header", "p.content", "p.image", "p.price", "p.author_id", "u.username", "p.created_at").
		From("posts p").
		Join("users u ON p.author_id = u.id").
		Where(squirrel.Eq{"p.id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build get post by ID query")
		return nil, fmt.Errorf("get post by ID query: %w", err)
	}
	var post entity.Post
	var username string
	err = a.db.QueryRow(ctx, query, args...).Scan(&post.ID, &post.Header, &post.Content, &post.Image, &post.Price, &post.AuthorID, &username, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		a.logger.WithError(err).Error("Failed to get post by ID")
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	post.AuthorUsername = username
	return &post, nil
}

func (a *PostAdapter) ListByAuthorID(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	queryBuilder := squirrel.Select("p.id", "p.header", "p.content", "p.image", "p.price", "p.author_id", "u.username", "p.created_at").
		From("posts p").
		Join("users u ON p.author_id = u.id").
		Where(squirrel.Eq{"p.author_id": authorID}).
		PlaceholderFormat(squirrel.Dollar)

	if minPrice, ok := filter["min_price"]; ok {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"p.price": minPrice})
	}
	if maxPrice, ok := filter["max_price"]; ok {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"p.price": maxPrice})
	}

	if sortBy == "" {
		sortBy = "created_at DESC"
	} else {
		validSortFields := map[string]bool{"created_at": true, "price": true}
		parts := strings.Split(sortBy, " ")
		if len(parts) != 2 || !validSortFields[parts[0]] || (parts[1] != "ASC" && parts[1] != "DESC") {
			return nil, 0, fmt.Errorf("invalid sortBy parameter")
		}
	}

	queryBuilder = queryBuilder.OrderBy(sortBy)

	// Запрос для подсчёта общего количества
	countQuery, countArgs, err := squirrel.Select("COUNT(*)").
		From("posts p").
		Where(squirrel.Eq{"p.author_id": authorID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build count query for posts by author")
		return nil, 0, fmt.Errorf("count query: %w", err)
	}
	var total int
	err = a.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, fmt.Errorf("no posts found for author")
		}
		a.logger.WithError(err).Error("Failed to count posts by author")
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	// Пагинация
	offset := (page - 1) * pageSize
	queryBuilder = queryBuilder.Limit(uint64(pageSize)).Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build list posts by author query")
		return nil, 0, fmt.Errorf("list posts query: %w", err)
	}

	rows, err := a.db.Query(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to list posts by author")
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		var username string
		err := rows.Scan(&post.ID, &post.Header, &post.Content, &post.Image, &post.Price, &post.AuthorID, &username, &post.CreatedAt)
		if err != nil {
			a.logger.WithError(err).Error("Failed to scan post row")
			return nil, 0, fmt.Errorf("scan post: %w", err)
		}
		post.AuthorUsername = username
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		a.logger.WithError(err).Error("Error iterating post rows")
		return nil, 0, fmt.Errorf("iterate posts: %w", err)
	}

	a.logger.WithFields(logrus.Fields{
		"author_id":   authorID,
		"page":        page,
		"page_size":   pageSize,
		"total_posts": total,
	}).Info("Posts by author listed from database")
	return posts, total, nil
}

func (a *PostAdapter) ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	queryBuilder := squirrel.Select("p.id", "p.header", "p.content", "p.image", "p.price", "p.author_id", "u.username", "p.created_at").
		From("posts p").
		Join("users u ON p.author_id = u.id").
		PlaceholderFormat(squirrel.Dollar)

	if minPrice, ok := filter["min_price"]; ok {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"p.price": minPrice})
	}
	if maxPrice, ok := filter["max_price"]; ok {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"p.price": maxPrice})
	}

	if sortBy == "" {
		sortBy = "created_at DESC"
	} else {
		validSortFields := map[string]bool{"created_at": true, "price": true}
		parts := strings.Split(sortBy, " ")
		if len(parts) != 2 || !validSortFields[parts[0]] || (parts[1] != "ASC" && parts[1] != "DESC") {
			return nil, 0, fmt.Errorf("invalid sortBy parameter")
		}
	}

	queryBuilder = queryBuilder.OrderBy(sortBy)

	// Запрос для подсчёта общего количества
	countQuery, countArgs, err := squirrel.Select("COUNT(*)").From("posts p").PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build count query for posts")
		return nil, 0, fmt.Errorf("count query: %w", err)
	}
	var total int
	err = a.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		a.logger.WithError(err).Error("Failed to count posts")
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	// Пагинация
	offset := (page - 1) * pageSize
	queryBuilder = queryBuilder.Limit(uint64(pageSize)).Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build list posts query")
		return nil, 0, fmt.Errorf("list posts query: %w", err)
	}

	rows, err := a.db.Query(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to list posts")
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		var username string
		err := rows.Scan(&post.ID, &post.Header, &post.Content, &post.Image, &post.Price, &post.AuthorID, &username, &post.CreatedAt)
		if err != nil {
			a.logger.WithError(err).Error("Failed to scan post row")
			return nil, 0, fmt.Errorf("scan post: %w", err)
		}
		post.AuthorUsername = username
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		a.logger.WithError(err).Error("Error iterating post rows")
		return nil, 0, fmt.Errorf("iterate posts: %w", err)
	}

	a.logger.WithFields(logrus.Fields{
		"page":        page,
		"page_size":   pageSize,
		"total_posts": total,
	}).Info("Posts listed from database")
	return posts, total, nil
}

func (a *PostAdapter) GetByHeaderAndContent(ctx context.Context, header, content string) (*entity.Post, error) {
	query, args, err := squirrel.Select("id", "header", "content", "image", "price", "author_id", "created_at").
		From("posts").
		Where(squirrel.Eq{"header": header, "content": content}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build get post by header and content query")
		return nil, fmt.Errorf("get post by header and content: %w", err)
	}
	var post entity.Post
	err = a.db.QueryRow(ctx, query, args...).Scan(&post.ID, &post.Header, &post.Content, &post.Image, &post.Price, &post.AuthorID, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		a.logger.WithError(err).Error("Failed to get post by header and content")
		return nil, fmt.Errorf("get post by header and content: %w", err)
	}
	return &post, nil
}

func (a *PostAdapter) Update(ctx context.Context, post *entity.Post) error {
	// Проверка на уникальность с исключением текущего поста
	queryCheck, argsCheck, err := squirrel.Select("id").
		From("posts").
		Where(squirrel.Eq{
			"header":    post.Header,
			"content":   post.Content,
			"author_id": post.AuthorID,
		}).
		Where(squirrel.NotEq{"id": post.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build check post query")
		return fmt.Errorf("check post existence: %w", err)
	}
	var existingID uuid.UUID
	err = a.db.QueryRow(ctx, queryCheck, argsCheck...).Scan(&existingID)
	if err == nil {
		a.logger.WithFields(logrus.Fields{
			"header":    post.Header,
			"author_id": post.AuthorID,
		}).Warn("Post with the same header, content, and author already exists")
		return fmt.Errorf("post with the same header, content, and author already exists")
	}
	if err != sql.ErrNoRows {
		a.logger.WithError(err).Error("Failed to check post existence")
		return fmt.Errorf("check post existence: %w", err)
	}

	// Обновление поста
	query, args, err := squirrel.Update("posts").
		Set("header", post.Header).
		Set("content", post.Content).
		Set("image", post.Image).
		Set("price", post.Price).
		Where(squirrel.Eq{"id": post.ID, "author_id": post.AuthorID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build update post query")
		return fmt.Errorf("update post query: %w", err)
	}
	result, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to update post")
		return fmt.Errorf("update post: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("post not found or not owned by user")
	}
	a.logger.WithFields(logrus.Fields{
		"post_id": post.ID,
	}).Info("Post updated in database")
	return nil
}

func (a *PostAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := squirrel.Delete("posts").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		a.logger.WithError(err).Error("Failed to build delete post query")
		return fmt.Errorf("delete post query: %w", err)
	}
	result, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		a.logger.WithError(err).Error("Failed to delete post")
		return fmt.Errorf("delete post: %w", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}
	a.logger.WithFields(logrus.Fields{
		"post_id": id,
	}).Info("Post deleted from database")
	return nil
}
