CREATE TABLE IF NOT EXISTS songs(
    id             SERIAL PRIMARY KEY,
    song           TEXT    NOT NULL,
    "release_date" DATE    NOT NULL,
    song_text      TEXT    NOT NULL,
    link           TEXT    NOT NULL,
    group_id       INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);