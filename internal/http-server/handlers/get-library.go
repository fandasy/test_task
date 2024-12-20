package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"test_task/internal/models"
	"test_task/internal/storage"
	"time"
)

// GetLibrary godoc
// @Summary Get library
// @Produce  json
// @Param offset query int false " "
// @Param limit query int false " "
// @Param group_id query int false " "
// @Param group query string false " "
// @Param song_id query int false " "
// @Param song query string false " "
// @Param release_date query string false " "
// @Param song_text query string false " "
// @Param link query string false " "
// @Success 200 {object} models.GetLibraryResponse
// @Success 404 {object} ErrResponse
// @Failure 400 {object} ErrResponse
// @Failure 500
// @Router /library [get]
func (h *Handler) GetLibrary(ctxTimeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		const fn = "handlers.GetLibrary"

		log := h.log.With(
			slog.String("fn", fn),
			slog.String("client_ip", c.ClientIP()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		offsetStr := c.Query("offset")
		limitStr := c.Query("limit")
		groupIDStr := c.Query("group_id")
		groupName := c.Query("group")
		songIDStr := c.Query("song_id")
		songName := c.Query("song")
		releaseDateStr := c.Query("release_date")
		songText := c.Query("song_text")
		link := c.Query("link")

		var (
			offset      int
			limit       int
			groupID     int
			songID      int
			releaseDate time.Time
			err         error
		)

		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				log.Debug("offset is not a number")

				c.JSON(http.StatusBadRequest, ErrResp("offset is not a number"))

				return
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				log.Debug("limit is not a number")

				c.JSON(http.StatusBadRequest, ErrResp("limit is not a number"))

				return
			}
		}

		if groupIDStr != "" {
			groupID, err = strconv.Atoi(groupIDStr)
			if err != nil {
				log.Debug("groupID is not a number")

				c.JSON(http.StatusBadRequest, ErrResp("groupID is not a number"))

				return
			}
		}

		if songIDStr != "" {
			songID, err = strconv.Atoi(songIDStr)
			if err != nil {
				log.Debug("songID is not a number")

				c.JSON(http.StatusBadRequest, ErrResp("songID is not a number"))

				return
			}
		}

		if releaseDateStr != "" {
			releaseDate, err = time.Parse("02.01.2006", releaseDateStr)
			if err != nil {
				log.Debug(err.Error())

				c.JSON(http.StatusBadRequest, ErrResp("release date is invalid"))

				return
			}
		}

		if limit < 0 {
			limit = 0
		}

		if offset < 0 {
			offset = 0
		}

		filters := &storage.GetLibraryFilters{
			Offset:      offset,
			Limit:       limit,
			GroupID:     groupID,
			GroupName:   groupName,
			SongID:      songID,
			SongName:    songName,
			ReleaseDate: releaseDate,
			SongText:    songText,
			Link:        link,
		}

		groupMap, err := h.db.GetLibrary(ctx, filters)
		if err != nil {
			if errors.Is(err, storage.ErrNothingFound) {
				log.Debug(err.Error(), slog.Any("filters", filters))

				c.JSON(http.StatusNotFound, ErrResp("nothing found"))

				return
			}
			log.Error(err.Error())

			c.Status(http.StatusInternalServerError)

			return
		}

		var response models.GetLibraryResponse

		for _, group := range groupMap {
			response.Library = append(response.Library, *group)
		}

		log.Debug("library data received successfully", slog.Any("filters", *filters))

		c.JSON(http.StatusOK, response)
	}
}
