package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"test_task/internal/storage"
	"time"
)

func (h *Handler) GetSongText(ctxTimeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		const fn = "handlers.GetSongText"

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

		offsetStr := c.Query("offset")
		limitStr := c.Query("limit")

		var (
			offset int
			limit  int
		)

		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				log.Debug("offset is not a number")

				c.JSON(http.StatusBadRequest, gin.H{"error": "offset is not a number"})

				return
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				log.Debug("limit is not a number")

				c.JSON(http.StatusBadRequest, gin.H{"error": "limit is not a number"})

				return
			}
		}

		song, err := h.db.GetSongText(ctx, int64(id))
		if err != nil {
			if errors.Is(err, storage.ErrSongNotFound) {
				log.Debug(err.Error())

				c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})

				return
			}
			log.Error(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		if offset == 0 && limit == 0 {
			log.Debug("lyrics sent unchanged")

			c.JSON(http.StatusOK, *song)

			return
		}

		song.SongText = splitIntoVerses(song.SongText, offset, limit)

		log.Debug("sent the filtered lyrics", slog.String("song_text", song.SongText))

		c.JSON(http.StatusOK, *song)
	}
}

func splitIntoVerses(text string, offset, limit int) string {
	textSlice := strings.Split(text, "\n")

	if offset >= len(textSlice) {
		return ""
	}

	if limit <= 0 {
		limit = len(textSlice)
	}

	if offset < 0 {
		offset = 0
	}

	end := offset + limit
	if end > len(textSlice) {
		end = len(textSlice)
	}

	fmt.Println(offset, end)

	return strings.Join(textSlice[offset:end], "\n")
}
