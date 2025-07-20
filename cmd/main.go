package main

import (
	"context"
	adapterPost "marketplace/internal/adapter/post"
	adapterUser "marketplace/internal/adapter/user"
	"marketplace/internal/handler"
	handlerAuth "marketplace/internal/handler/auth"
	handlerPost "marketplace/internal/handler/post"
	handlerUser "marketplace/internal/handler/user"
	serviceAuth "marketplace/internal/service/auth"
	servicePost "marketplace/internal/service/post"
	serviceUser "marketplace/internal/service/user"
	usecaseAuth "marketplace/internal/usecase/auth"
	usecasePost "marketplace/internal/usecase/post"
	usecaseUser "marketplace/internal/usecase/user"
	"marketplace/pkg/config"
	"marketplace/pkg/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Настройка логгера
	log := logger.SetupLogger(cfg.Logger.Level, cfg.Logger.Format)

	// Подключение к базе данных
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer dbPool.Close()

	// Инициализация адаптеров
	postAdapter := adapterPost.NewPostAdapter(dbPool, log)
	userAdapter := adapterUser.NewUserAdaper(dbPool, log)

	// Инициализация AuthService
	authImpl := usecaseAuth.NewAuthImpl(cfg.JWT.SecretKey)

	// Инициализация usecases
	userUsecase := usecaseUser.NewUserUseCase(userAdapter, authImpl, log)
	postUsecase := usecasePost.NewPostUsecase(postAdapter, userAdapter, authImpl, log)

	// Инициализация сервисов
	authService := serviceAuth.NewAuthService(authImpl, log)
	userService := serviceUser.NewUserService(userUsecase, log)
	postService := servicePost.NewPostService(postUsecase, log)

	// Инициализация обработчиков
	authHandler := handlerAuth.NewAuthHandler(authService, log)
	userHandler := handlerUser.NewUserHandler(userService, log)
	postHandler := handlerPost.NewPostHandler(postService, userService, log)

	// Настройка маршрутов
	router := handler.NewRouter(userHandler, postHandler, authHandler)
	ginRouter := router.SetupRoutes()

	// Запуск сервера
	log.Infof("Starting server on port %s", cfg.Server.Port)
	if err := ginRouter.Run(cfg.Server.Port); err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
