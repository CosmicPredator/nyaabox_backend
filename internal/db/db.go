package db

import (
	"cosmic/nyaabox/internal/config"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2/log"
	_ "modernc.org/sqlite"
)

type FileEntry struct {
	FileId    string
	FileName  string
	FilePath  string
	ExpiresAt time.Time
}

type DbHandler struct {
	handle *sql.DB
}

func (db *DbHandler) Init() error {
	log.Debug("executing DB init query")
	_, err := db.handle.Exec(createTableQuery)
	return err
}

func (db *DbHandler) AddEntry(fileEntry *FileEntry) error {
	log.Debugf("executing DB AddEntry query for: %v", fileEntry)
	_, err := db.handle.Exec(
		insertEntryQuery,
		fileEntry.FileId,
		fileEntry.FileName,
		fileEntry.FilePath,
		fileEntry.ExpiresAt,
	)
	return err
}

func (db *DbHandler) GetEntryByFileId(fileId string) (*FileEntry, error) {
	log.Debugf("executing DB GetEntry query for: %v", fileId)
	var fileEntry FileEntry
	err := db.handle.QueryRow(getEntryQuery, fileId).
		Scan(&fileEntry.FileId, &fileEntry.FileName, &fileEntry.FilePath, &fileEntry.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &fileEntry, nil
}

func (db *DbHandler) GetEntriesByExpiredAt(expiresAt time.Time) ([]*FileEntry, error) {
	log.Debugf("executing DB GetEntriesByExpiredAt query for stamp: %v", expiresAt)
	var fileEntries []*FileEntry
	rows, err := db.handle.Query(getExpiresAtQuery, expiresAt)
	if err != nil { return nil, err }
	defer rows.Close()
	
	for rows.Next() {
		f := new(FileEntry)
		if err != rows.Scan(&f.FileId, &f.FileName, &f.FilePath) {
			return nil, err
		}
		fileEntries = append(fileEntries, f)
	}
	
	return fileEntries, nil
}

func (db *DbHandler) DeleteEntriesByExpiredAt(expiresAt time.Time) error {
	res, err := db.handle.Exec(deleteExpiredQuery, expiresAt)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	log.Debugf("deleted %d expired rows from DB", count)
	return nil
}

func (db *DbHandler) Close() {
	log.Debugf("trying to close DB handle")
	db.handle.Close()
}

func NewDbHandler() (*DbHandler, error) {
	var dbPath string
	dbPath, err := config.GetEnv("NB_DB_PATH")
	if err != nil {
		log.Warn("'NB_DB_PATH' not provided. defaulting to './nyaabox.db'")
		dbPath = "nyaabox.db"
	}

	log.Debug("trying to initialize DB")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	log.Info("DB initialized successfully")
	return &DbHandler{handle: db}, nil
}
