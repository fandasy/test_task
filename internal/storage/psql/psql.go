package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
	"strings"
	"test_task/internal/config"
	"test_task/internal/storage"
	"test_task/pkg/e"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config, migratePath string) (*Storage, error) {
	const fn = "psql.New"

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, e.Wrap(fn, err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap(fn, err)
	}

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, e.Wrap(fn, err)
	}

	m, err := migrate.NewWithDatabaseInstance(migratePath, "postgres", migrationDriver)
	if err != nil {
		return nil, e.Wrap(fn, err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info(fn, slog.String("msg", "Migrations: no change to apply"))
		} else {
			return nil, e.Wrap(fn, err)
		}
	} else {
		slog.Info(fn, slog.String("msg", "Migrations applied successfully!"))
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveGroup(ctx context.Context, groupName string) (int64, error) {
	const fn = "psql.SaveGroup"

	q := `
	INSERT INTO groups (group_name)
	VALUES ($1)
	RETURNING id;`

	var groupID int64

	if err := s.db.QueryRowContext(ctx, q, groupName).Scan(&groupID); err != nil {
		return 0, e.Wrap(fn, err)
	}

	return groupID, nil
}

func (s *Storage) SaveSong(ctx context.Context, songInfo *storage.SongInfo) (int64, error) {
	const fn = "psql.SaveSong"

	q := `
	INSERT INTO songs (song, release_date, song_text, link, group_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`

	args := []any{
		songInfo.Song,
		songInfo.Date,
		songInfo.Text,
		songInfo.Link,
		songInfo.GroupID,
	}

	var songID int64

	if err := s.db.QueryRowContext(ctx, q, args...).Scan(&songID); err != nil {
		return 0, e.Wrap(fn, err)
	}

	return songID, nil
}

func (s *Storage) GroupExists(ctx context.Context, GroupName string) (int64, bool, error) {
	const fn = "psql.GroupExists"

	q := `SELECT id FROM groups WHERE group_name = $1;`

	var groupID int64

	if err := s.db.QueryRowContext(ctx, q, GroupName).Scan(&groupID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}

		return 0, false, e.Wrap(fn, err)
	}

	return groupID, true, nil
}

func (s *Storage) SongExists(ctx context.Context, SongName string, GroupID int64) (int64, bool, error) {
	const fn = "psql.SongExists"

	q := `SELECT id FROM songs WHERE song = $1 AND group_id = $2;`

	var songID int64

	if err := s.db.QueryRowContext(ctx, q, SongName, GroupID).Scan(&songID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}

		return 0, false, e.Wrap(fn, err)
	}

	return songID, true, nil
}

func (s *Storage) DeleteSong(ctx context.Context, songID int) error {
	const fn = "psql.DeleteSong"

	q := `DELETE FROM songs WHERE id = $1;`

	res, err := s.db.ExecContext(ctx, q, songID)
	if err != nil {
		return e.Wrap(fn, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return e.Wrap(fn, err)
	}

	if rowsAffected == 0 {
		return e.Wrap(fn, storage.ErrSongNotFound)
	}

	return nil
}

func (s *Storage) GetSongText(ctx context.Context, songID int64) (*storage.SongResp, error) {
	const fn = "psql.GetSongText"

	var songResp storage.SongResp

	q := `SELECT song, song_text FROM songs WHERE id = $1;`

	if err := s.db.QueryRowContext(ctx, q, songID).Scan(&songResp.SongName, &songResp.SongText); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrSongNotFound
		}

		return nil, e.Wrap(fn, err)
	}

	songResp.SongID = songID

	return &songResp, nil
}

func (s *Storage) GetLibrary(ctx context.Context, filters *storage.GetLibraryFilters) (map[int64]*storage.GroupResp, error) {
	const fn = "psql.GetLibrary"

	query := `
	SELECT g.id, g.group_name, s.id, s.song, s.release_date, s.song_text, s.link 
	FROM groups g
	LEFT JOIN songs s ON g.id = s.group_id
	`

	var args []interface{}
	var sets []string
	paramIndex := 1

	if filters.GroupName != "" {
		sets = append(sets, fmt.Sprintf("g.group_name = $%d", paramIndex))
		args = append(args, filters.GroupName)
		paramIndex++
	}
	if filters.GroupID != 0 {
		sets = append(sets, fmt.Sprintf("g.id = $%d", paramIndex))
		args = append(args, filters.GroupID)
		paramIndex++
	}
	if filters.SongName != "" {
		sets = append(sets, fmt.Sprintf("s.song = $%d", paramIndex))
		args = append(args, filters.SongName)
		paramIndex++
	}
	if filters.SongID != 0 {
		sets = append(sets, fmt.Sprintf("s.id = $%d", paramIndex))
		args = append(args, filters.SongID)
		paramIndex++
	}
	if !filters.ReleaseDate.IsZero() {
		sets = append(sets, fmt.Sprintf("s.release_date = $%d", paramIndex))
		args = append(args, filters.ReleaseDate)
		paramIndex++
	}
	if filters.SongText != "" {
		sets = append(sets, fmt.Sprintf("s.song_text = $%d", paramIndex))
		args = append(args, filters.SongText)
		paramIndex++
	}
	if filters.Link != "" {
		sets = append(sets, fmt.Sprintf("s.link = $%d", paramIndex))
		args = append(args, filters.Link)
		paramIndex++
	}

	if len(sets) > 0 {
		query += "WHERE "
		query += strings.Join(sets, "AND ")
	}

	query += fmt.Sprintf(" ORDER BY g.id")

	if filters.Offset != 0 {
		query += fmt.Sprintf(" OFFSET $%d", paramIndex)
		args = append(args, filters.Offset)
		paramIndex++
	}

	if filters.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramIndex)
		args = append(args, filters.Limit)
		paramIndex++
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, e.Wrap(fn, err)
	}
	defer rows.Close()

	groupMap := make(map[int64]*storage.GroupResp)

	for rows.Next() {
		var (
			g  storage.GroupResp
			s  storage.SongResp
			rd time.Time
		)

		err := rows.Scan(&g.GroupID, &g.GroupName, &s.SongID, &s.SongName, &rd, &s.SongText, &s.Link)
		if err != nil {
			return nil, e.Wrap(fn, storage.ErrNothingFound)
		}

		s.ReleaseDate = rd.Format("02.01.2006")

		if _, exists := groupMap[g.GroupID]; !exists {
			groupMap[g.GroupID] = &storage.GroupResp{
				GroupID:   g.GroupID,
				GroupName: g.GroupName,
				SongInfo:  []storage.SongResp{},
			}
		}

		if s.SongID != 0 {
			groupMap[g.GroupID].SongInfo = append(groupMap[g.GroupID].SongInfo, s)
		}
	}

	if len(groupMap) == 0 {
		return nil, e.Wrap(fn, storage.ErrNothingFound)
	}

	return groupMap, nil
}

func (s *Storage) UpdateSong(ctx context.Context, songID int, songInfo *storage.SongInfo) error {
	const fn = "psql.UpdateSong"

	query := "UPDATE songs SET "
	var args []interface{}
	var sets []string
	paramIndex := 1

	if songInfo.Song != "" {
		sets = append(sets, fmt.Sprintf("song = $%d", paramIndex))
		args = append(args, songInfo.Song)
		paramIndex++
	}
	if !songInfo.Date.IsZero() {
		sets = append(sets, fmt.Sprintf("release_date = $%d", paramIndex))
		args = append(args, songInfo.Date)
		paramIndex++
	}
	if songInfo.Text != "" {
		sets = append(sets, fmt.Sprintf("song_text = $%d", paramIndex))
		args = append(args, songInfo.Text)
		paramIndex++
	}
	if songInfo.Link != "" {
		sets = append(sets, fmt.Sprintf("link = $%d", paramIndex))
		args = append(args, songInfo.Link)
		paramIndex++
	}

	if len(sets) == 0 {
		return e.Wrap(fn, storage.ErrNoFieldsUpdate)
	}

	query += strings.Join(sets, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", paramIndex)
	args = append(args, songID)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return e.Wrap(fn, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return e.Wrap(fn, err)
	}

	if rowsAffected == 0 {
		return e.Wrap(fn, storage.ErrSongNotFound)
	}

	return nil
}
