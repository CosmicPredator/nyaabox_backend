package db

const createTableQuery = `
CREATE TABLE IF NOT EXISTS files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	file_id TEXT NOT NULL UNIQUE,
	filename TEXT NOT NULL,
	filepath TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL
);`

const insertEntryQuery = `
INSERT INTO files (file_id, filename, filepath, expires_at)
VALUES (?, ?, ?, ?);`

const getEntryQuery = `
SELECT file_id, filename, filepath, expires_at FROM files 
WHERE file_id = ?`

const getExpiresAtQuery = `
SELECT file_id, filename, filepath FROM files WHERE expires_at <= ?`

const deleteExpiredQuery = `
DELETE FROM files WHERE expires_at <= ?`