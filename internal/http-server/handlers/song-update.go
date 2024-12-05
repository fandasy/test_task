package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"test_task/internal/lib/l/sl"
	"test_task/internal/models"
	"test_task/internal/storage"
	"test_task/pkg/validate"
	"time"
)

type SongUpdateRequest struct {
	SongName    *string `json:"song_name,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
	SongText    *string `json:"song_text,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type parsedSongUpdateReq struct {
	songName    string
	releaseDate string
	songText    string
	link        string
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

			c.JSON(http.StatusBadRequest, ErrResp("id is empty"))

			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Debug("id is invalid", slog.Int("ID", id))

			c.JSON(http.StatusBadRequest, ErrResp("id is invalid"))

			return
		}

		log.Debug("id is valid", slog.Int("songID", id))

		var req SongUpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request", sl.Err(err))

			c.JSON(http.StatusBadRequest, ErrResp("incorrect json"))

			return
		}

		preq := parsingReq(req)

		log.Debug("request body decoded", slog.Any("req", preq))

		var releaseDate time.Time

		if preq.releaseDate != "" {
			releaseDate, err = time.Parse("02.01.2006", preq.releaseDate)
			if err != nil {
				log.Debug(err.Error())

				c.JSON(http.StatusBadRequest, ErrResp("release date is invalid, correct format: DD.MM.YYYY"))

				return
			}
		}

		if !(preq.link == "" || preq.link == "NULL") {
			if !validate.Link(ctx, preq.link) {
				log.Debug("request link is invalid", slog.String("link", preq.link))

				c.JSON(http.StatusBadRequest, ErrResp("link is invalid"))

				return
			}
		}

		songInfo := &storage.SongInfo{
			Song: preq.songName,
			Date: releaseDate,
			Text: preq.songText,
			Link: preq.link,
		}

		if err := h.db.UpdateSong(ctx, id, songInfo); err != nil {
			switch {
			case errors.Is(err, storage.ErrNoFieldsUpdate):
				log.Debug(err.Error())

				c.JSON(http.StatusBadRequest, ErrResp("no fields to update"))

				return
			case errors.Is(err, storage.ErrSongNotFound):
				log.Debug(err.Error())

				c.JSON(http.StatusNotFound, ErrResp("song not found"))

				return
			default:
				log.Debug(err.Error())

				c.Status(http.StatusInternalServerError)

				return
			}
		}

		log.Debug("song data update",
			slog.Int("songID", id),
			slog.Any("data", preq),
		)

		c.JSON(http.StatusOK, models.SongUpdateResponse{
			SongID: id,
			UpdateInfo: models.UpdateInfo{
				SongName:    preq.songName,
				ReleaseDate: preq.releaseDate,
				SongText:    preq.songText,
				Link:        preq.link,
			}})
	}
}

func parsingReq(input interface{}) parsedSongUpdateReq {
	var pr parsedSongUpdateReq
	val := reflect.ValueOf(input)

	numFields := val.NumField()

	for i := 0; i < numFields; i++ {
		value := val.Field(i)

		if value.Kind() == reflect.Ptr && !value.IsNil() {
			elemValue := value.Elem().String()
			switch i {
			case 0:
				if elemValue == "" {
					pr.songName = "NULL"
				} else {
					pr.songName = elemValue
				}
			case 1:
				if elemValue == "" {
					pr.releaseDate = "01.01.0001"
				} else {
					pr.releaseDate = elemValue
				}
			case 2:
				if elemValue == "" {
					pr.songText = "NULL"
				} else {
					pr.songText = elemValue
				}
			case 3:
				if elemValue == "" {
					pr.link = "NULL"
				} else {
					pr.link = elemValue
				}
			}
		}
	}

	return pr
}
