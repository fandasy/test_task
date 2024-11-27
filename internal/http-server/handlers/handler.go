package handlers

import (
	"log/slog"
	your_api "test_task/internal/clients/your-api"
	"test_task/internal/storage"
)

type Handler struct {
	db      storage.Storage
	log     *slog.Logger
	yourApi *your_api.Client
}

func New(db storage.Storage, log *slog.Logger, yourApi *your_api.Client) *Handler {
	return &Handler{
		db:      db,
		log:     log,
		yourApi: yourApi,
	}
}
