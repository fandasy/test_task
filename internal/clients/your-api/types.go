package your_api

import "errors"

type Response struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

var (
	ErrBadRequest          = errors.New("bad Request")
	ErrInternalServerError = errors.New("internal Server Error")
)
