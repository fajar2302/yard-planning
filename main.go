package main

import (
	"yard-planning/internal/config"
	"yard-planning/internal/db"
	"yard-planning/internal/handler"
	"yard-planning/internal/repository"
	"yard-planning/internal/routes"
	"yard-planning/internal/services"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, _ := config.Load()
	db, _ := db.NewPostgres(cfg)
	repo := repository.NewRepo(db)
	service := services.NewService(repo)
	handler := handler.NewHandler(service)
	e := echo.New()
	routes.RegisterRoutes(e, handler)
	e.Start(":8080")
}
