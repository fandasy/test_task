package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"errors"
	your_api "test_task/internal/clients/your-api"
	"test_task/internal/lib/l/sl"
	"test_task/internal/models"
	"test_task/internal/storage"
	"time"
)

type SaveSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

// SaveSong godoc
// @Summary Save song
// @Accept json
// @Produce json
// @Param song body SaveSongRequest true "Group and Song name"
// @Success 200 {object} models.SaveSongResponse
// @Success 404 {object} ErrResponse
// @Failure 400 {object} ErrResponse
// @Failure 500
// @Router /song [post]
func (h *Handler) SaveSong(ctxTimeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		const fn = "handlers.SaveSong"

		log := h.log.With(
			slog.String("fn", fn),
			slog.String("client_ip", c.ClientIP()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		var req SaveSongRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			c.JSON(http.StatusBadRequest, ErrResp("incorrect json"))

			return
		}

		log.Debug("request body decoded", slog.Any("req", req))

		groupID, groupExists, err := h.db.GroupExists(ctx, req.Group)
		if err != nil {
			log.Error("failed to check if group exists", sl.Err(err))

			c.Status(http.StatusInternalServerError)

			return
		}

		if groupExists {
			log.Debug("group already exists", slog.Int64("groupID", groupID))

			songID, songExists, err := h.db.SongExists(ctx, req.Song, groupID)
			if err != nil {
				log.Error("failed to check if song exists", sl.Err(err))

				c.Status(http.StatusInternalServerError)

				return
			}

			if songExists {
				log.Debug("existing data sent",
					slog.Int64(req.Group, groupID),
					slog.Int64(req.Song, songID))

				c.JSON(http.StatusOK, models.SaveSongResponse{
					GroupID: groupID,
					SongID:  songID,
				})

				return
			}
		}

		resp, err := h.yourApi.GetSongInfo(ctx, req.Group, req.Song)
		if err != nil {
			log.Error("failed to get song info", sl.Err(err))

			switch {
			case errors.Is(err, context.DeadlineExceeded):
				c.JSON(http.StatusRequestTimeout, ErrResp("request took too long"))

			case errors.Is(err, your_api.ErrBadRequest):
				c.JSON(http.StatusBadRequest, ErrResp("bad request"))

			default:
				c.Status(http.StatusInternalServerError)
			}

			return
		}

		log.Debug("your api response body decoded", slog.Any("resp", resp))

		if !groupExists {
			groupID, err = h.db.SaveGroup(ctx, req.Group)
			if err != nil {
				log.Error("failed to save group", sl.Err(err))

				c.Status(http.StatusInternalServerError)

				return
			}

			log.Debug("group saved", slog.Int64(req.Group, groupID))
		}

		releaseDate, _ := time.Parse("02.01.2006", resp.ReleaseDate)

		songInfo := &storage.SongInfo{
			Song:    req.Song,
			Date:    releaseDate,
			Text:    resp.Text,
			Link:    resp.Link,
			GroupID: groupID,
		}

		songID, err := h.db.SaveSong(ctx, songInfo)
		if err != nil {
			log.Error("failed to save song", sl.Err(err))

			c.Status(http.StatusInternalServerError)

			return
		}

		log.Debug("data save",
			slog.Int64(req.Group, groupID),
			slog.Int64(req.Song, songID))

		c.JSON(http.StatusOK, models.SaveSongResponse{
			GroupID: groupID,
			SongID:  songID,
		})
	}
}
