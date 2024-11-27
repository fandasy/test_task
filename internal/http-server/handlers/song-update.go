package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"test_task/internal/lib/l/sl"
	"test_task/internal/models"
	"test_task/internal/storage"
	"test_task/pkg/validate"
	"time"
)

type SongUpdateRequest struct {
	SongName    string `json:"song_name,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	SongText    string `json:"song_text,omitempty"`
	Link        string `json:"link,omitempty"`
}

// SongUpdate godoc
// @Summary Update song data
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param update_data body SongUpdateRequest true "Update data"
// @Success 200 {object} models.SongUpdateResponse
// @Success 404 {object} ErrResponse
// @Failure 400 {object} ErrResponse
// @Failure 500
// @Router /song/{id} [patch]
func (h *Handler) SongUpdate(ctxTimeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		const fn = "handlers.SongUpdate"

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

		var req SongUpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect json"})

			return
		}

		log.Debug("request body decoded", slog.Any("req", req))

		var releaseDate time.Time

		if req.ReleaseDate != "" {
			releaseDate, err = time.Parse("02.01.2006", req.ReleaseDate)
			if err != nil {
				log.Debug(err.Error())

				c.JSON(http.StatusBadRequest, gin.H{"error": "release date is invalid, correct format: 16.07.2006"})

				return
			}
		}

		if req.Link != "" {
			if !validate.Link(ctx, req.Link) {
				log.Debug("request link is invalid", slog.String("link", req.Link))

				c.JSON(http.StatusBadRequest, gin.H{"error": "link is invalid"})

				return
			}
		}

		songInfo := &storage.SongInfo{
			Song: req.SongName,
			Date: releaseDate,
			Text: req.SongText,
			Link: req.Link,
		}

		if err := h.db.UpdateSong(ctx, id, songInfo); err != nil {
			switch {
			case errors.Is(err, storage.ErrNoFieldsUpdate):
				log.Debug(err.Error())

				c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})

				return
			case errors.Is(err, storage.ErrSongNotFound):
				log.Debug(err.Error())

				c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})

				return
			default:
				log.Debug(err.Error())

				c.Status(http.StatusInternalServerError)

				return
			}
		}

		log.Debug("song data update",
			slog.Int("songID", id),
			slog.Any("data", req),
		)

		c.JSON(http.StatusOK, models.SongUpdateResponse{
			SongID: id,
			UpdateInfo: models.UpdateInfo{
				SongName:    req.SongName,
				ReleaseDate: req.ReleaseDate,
				SongText:    req.SongText,
				Link:        req.Link,
			}})
	}
}
