package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/entity"
	usecaseAuth "marketplace/internal/usecase/auth"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	userRepo UserRepository
	authRepo usecaseAuth.AuthService
	logger   *logrus.Logger
}

func NewUserUseCase(userRepo UserRepository, authRepo usecaseAuth.AuthService, logger *logrus.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		authRepo: authRepo,
		logger:   logger,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, "", fmt.Errorf("validate password: %w", err)
	}
	if _, err := uc.userRepo.GetByUsername(ctx, username); err == nil {
		return nil, "", fmt.Errorf("username already exists")
	}

	hashedPassword, err := uc.authRepo.GeneratePasswordHash(password)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		ID:             uuid.New(),
		Username:       username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, "", fmt.Errorf("validate user: %w", err)
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}

	token, err := uc.authRepo.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("generate jwt: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.ID,
	}).Info("User registered")

	return user.ToDTO(), token, nil
}

func (uc *UserUseCase) Login(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", fmt.Errorf("get user: %w", err)
	}

	if err := uc.authRepo.VerifyPassword(user.HashedPassword, password); err != nil {
		return nil, "", fmt.Errorf("verify password: %w", err)
	}

	token, err := uc.authRepo.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("generate jwt: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.ID,
	}).Info("User logged in")

	return user.ToDTO(), token, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User fetched")

	return user, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id uuid.UUID, username, password string) error {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok || userID != id {
		return fmt.Errorf("unauthorized")
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if username != "" && username != user.Username {
		_, err := uc.userRepo.GetByUsername(ctx, username)
		if err == nil {
			return fmt.Errorf("username %s already exists", username)
		}
		user.Username = username
	}

	if password != "" {
		hashedPassword, err := uc.authRepo.GeneratePasswordHash(password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		user.HashedPassword = hashedPassword
	}

	if username == "" && password == "" {
		return fmt.Errorf("nothing to update")
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User updated")

	return nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok || userID != id {
		return fmt.Errorf("unauthorized")
	}

	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted")

	return nil
}
