package storage

import (
	"cosmic/nyaabox/internal/config"
	"cosmic/nyaabox/internal/db"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func SaveToDisk(ctx *fiber.Ctx, dbHandler db.DbHandler, file *multipart.FileHeader) (string, error) {
	log.Debugf("Processing: %s | %s | %s", file.Filename, file.Size, file.Header["Content-Type"][0])
	fileExt := filepath.Ext(file.Filename)
	fileId := generateFileId(6)
	fileName := fileId + fileExt
	filePath := filepath.Join(GetSaveDir(), fileName)
	
	fileStream, err := file.Open()
	if err != nil { return "", err }
	defer fileStream.Close()
	
	outFile, err := os.Create(filePath)
	if err != nil { return "", err }
	defer outFile.Close()
	
	_, err = io.Copy(outFile, fileStream)
	if err != nil { return "", err }
	
	expireDays, err := config.GetEnv("NB_FILE_RETENTION_DAYS")
	if err != nil {
		expireDays = "3"
	}
	expireDaysInt, _ := strconv.Atoi(expireDays)
	
	err = dbHandler.AddEntry(&db.FileEntry{
		FileId: fileId,
		FileName: fileName,
		FilePath: filePath,
		ExpiresAt: time.Now().Add(time.Duration(expireDaysInt) * 24 * time.Hour),
	})
	
	if err != nil { return "", err }
	return fileName, nil
}

func GetSaveDir() string {
	savePath, err := config.GetEnv("NB_SAVE_PATH")
	if err != nil {
		log.Warn("'NB_SAVE_PATH' not provided. defaulting to '.'")
		savePath = "."
	}
	return savePath
}

func PrepareSaveDir() (string, error) {
	savePath, err := config.GetEnv("NB_SAVE_PATH")
	
	if err != nil {
		log.Warn("'NB_SAVE_PATH' not provided. defaulting to '.'")
		savePath = "."
	}
	
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		log.Infof("directory %s does not exist. creating one", savePath)
		err := os.MkdirAll(savePath, 0755)
		if err != nil {
			return "", err
		}
	}
	
	return savePath, nil
}

func GetFile(fileId string, dbHandle *db.DbHandler) (*db.FileEntry, error) {
	fileInfo, err := dbHandle.GetEntryByFileId(
		removeFileExtension(fileId),
	)
	if err != nil {
		return nil, fmt.Errorf("file expired or not found: %s", fileId)
	}
	return fileInfo, nil
}

func RemoveFile(filePath string) (err error) {
	err = os.Remove(filePath)
	return
}

func removeFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

func generateFileId(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}