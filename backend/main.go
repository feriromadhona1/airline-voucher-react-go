package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"airline-voucher/internal/database"
	"airline-voucher/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(".", "vouchers.db")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	store, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer store.Close()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	voucherHandler := handlers.NewVoucherHandler(store)
	api := e.Group("/api")
	api.POST("/check", voucherHandler.Check)
	api.POST("/generate", voucherHandler.Generate)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Printf("backend listening on :%s (db=%s)", port, dbPath)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
