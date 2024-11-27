package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"test_task/internal/storage"
	"time"
)

func (h *Handler) DeleteSong(ctxTimeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		const fn = "handlers.DeleteSong"

		log := h.log.With(
			slog.String("fn", fn),
			slog.String("client_ip", c.ClientIP()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		idStr := c.Param("id")
		if idStr == "" {
			log.Debug("id is empty")

			c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})

			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Debug("id is invalid", slog.Int("ID", id))

			c.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid"})

			return
		}

		log.Debug("id is valid", slog.Int("songID", id))

		if err := h.db.DeleteSong(ctx, id); err != nil {
			if errors.Is(err, storage.ErrSongNotFound) {
				log.Debug(err.Error())

				c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})

				return
			}

			log.Error(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

			return
		}

		log.Debug("song deleted", slog.Int("songID", id))

		c.JSON(http.StatusOK, gin.H{"message": "song deleted", "songID": id})
	}
}
