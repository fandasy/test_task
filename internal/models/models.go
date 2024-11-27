package models

type Song struct {
	SongID      int64  `json:"song_id"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	SongText    string `json:"song_text"`
	Link        string `json:"link"`
}

type Group struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	SongInfo  []Song `json:"song_info"`
}

type DeleteSongResp struct {
	Message string `json:"message"`
	SongID  int    `json:"song_id"`
}

type GetLibraryResponse struct {
	Library []Group `json:"library"`
}

type SongTextResp struct {
	SongID   int64  `json:"song_id"`
	SongName string `json:"song_name"`
	SongText string `json:"song_text"`
}

type SaveSongResponse struct {
	GroupID int64 `json:"group_Id"`
	SongID  int64 `json:"song_Id"`
}

type SongUpdateResponse struct {
	SongID     int        `json:"song_id"`
	UpdateInfo UpdateInfo `json:"update_info"`
}

type UpdateInfo struct {
	SongName    string `json:"song_name,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	SongText    string `json:"song_text,omitempty"`
	Link        string `json:"link,omitempty"`
}
