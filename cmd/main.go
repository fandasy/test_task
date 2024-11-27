package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	_ "test_task/docs"
	your_api "test_task/internal/clients/your-api"
	"test_task/internal/config"
	"test_task/internal/http-server/handlers"
	"test_task/internal/http-server/middleware/cors"
	"test_task/internal/http-server/middleware/logger"
	"test_task/internal/lib/l"
	"test_task/internal/storage/psql"
	"test_task/pkg/e"
	"time"
)

// @title           Music Library API
// @version         1.0.0
// @description     API for managing a music library
func main() {

	cfg, err := config.LoadEnvConfig("config.env")
	if err != nil {
		panic(err)
	}

	log := l.SetupLogger(cfg.Slog)

	storage, err := psql.New(cfg, "file://migrations")
	if err != nil {
		panic(err)
	}

	yourApiClient := your_api.NewClient(cfg.YourAPIHost)

	handler := handlers.New(storage, log, yourApiClient)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(cors.Middleware()) // for OPTIONS requests
	router.Use(logger.Middleware(log))
	router.Use(gin.Recovery())

	router.POST("/song", handler.SaveSong(30*time.Second))
	router.GET("/library", handler.GetLibrary(30*time.Second))
	router.GET("/song/:id/text", handler.GetSongText(30*time.Second))
	router.DELETE("/song/:id", handler.DeleteSong(30*time.Second))
	router.PATCH("/song/:id", handler.SongUpdate(30*time.Second))

	router.GET("/swagger/:any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info("server starting", slog.String("address", cfg.Addr))

	srv := &http.Server{
		Addr:        cfg.Addr,
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(e.Wrap("failed to start server", err))
		}
	}()

	// Server shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop
	log.Info("got signal", slog.String("signal", sign.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info("failed to shutdown server", slog.String("error", err.Error()))
	}

	log.Info("server shutdown", slog.String("address", srv.Addr))
}
