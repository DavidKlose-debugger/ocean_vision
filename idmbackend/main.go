package main

import (
	"log"

	"github.com/MoyInGxing/idm/app"
	"github.com/MoyInGxing/idm/config"
	"github.com/MoyInGxing/idm/handler"
	"github.com/MoyInGxing/idm/infra/database"
	"github.com/MoyInGxing/idm/internal/myrouter"
	"github.com/MoyInGxing/idm/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	db := database.GetDB()

	userRepo := database.NewGORMUserRepository(db)
	sessionRepo := database.NewGORMSessionRepository(db)
	speciesRepo := database.NewGORMSpeciesRepository(db)

	userService := app.NewUserService(userRepo)
	authService := app.NewAuthService(userRepo, sessionRepo, cfg)
	speciesService := app.NewSpeciesService(speciesRepo)

	userHandler := handler.NewUserHandler(userService, authService)
	speciesHandler := handler.NewSpeciesHandler(speciesService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(authService)

	r := myrouter.SetupRouter(userHandler, speciesHandler, authMiddleware, adminAuthMiddleware)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
