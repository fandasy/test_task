package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	SaveGroup(ctx context.Context, groupName string) (int64, error)
	SaveSong(ctx context.Context, songInfo *SongInfo) (int64, error)
	GroupExists(ctx context.Context, GroupName string) (int64, bool, error)
	SongExists(ctx context.Context, SongName string, GroupID int64) (int64, bool, error)
	DeleteSong(ctx context.Context, songID int) error
	GetSongText(ctx context.Context, songID int64) (*SongResp, error)
	GetLibrary(ctx context.Context, filters *GetLibraryFilters) (map[int64]*GroupResp, error)
	UpdateSong(ctx context.Context, songID int, songInfo *SongInfo) error
}

var (
	ErrSongNotFound   = errors.New("song not found")
	ErrNoFieldsUpdate = errors.New("no fields to update")
	ErrNothingFound   = errors.New("nothing found")
)

type SongInfo struct {
	SongID  int
	Song    string
	Date    time.Time
	Text    string
	Link    string
	GroupID int64
}

type GetLibraryFilters struct {
	Offset      int
	Limit       int
	GroupID     int
	GroupName   string
	SongID      int
	SongName    string
	ReleaseDate time.Time
	SongText    string
	Link        string
}

type SongResp struct {
	SongID      int64  `json:"song_id"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date,omitempty"`
	SongText    string `json:"song_text"`
	Link        string `json:"link,omitempty"`
}

type GroupResp struct {
	GroupID   int64      `json:"group_id"`
	GroupName string     `json:"group_name"`
	SongInfo  []SongResp `json:"song_info"`
}
