package main

import (
	"code/internal/config"
	"code/internal/handlers"
	"code/templates"
	"html/template"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	router := gin.Default()

	tmpl := template.Must(template.ParseFS(templates.FS, "**/*.html"))
	router.SetHTMLTemplate(tmpl)

	userHandler := handlers.NewUserHandler(db)

	router.GET("/", userHandler.Index)
	router.GET("/users/list", userHandler.List)
	router.POST("/users", userHandler.Create)
	router.GET("/users/:id/edit", userHandler.EditForm)
	router.PUT("/users/:id", userHandler.Update)
	router.GET("/users/edit/clear", userHandler.ClearEditForm)
	router.DELETE("/users/:id", userHandler.Delete)

	slog.Info("Server is running on port 8080")
	if err := router.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
