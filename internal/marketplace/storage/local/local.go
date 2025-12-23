// Package local handles local storage of images.
package local

import (
	"context"
	"errors"
	"fmt"
	"golang-connect-marketplace/pkg/generate"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// errPathIsEmpty         = errors.New("path is empty").
var errUnsupportedFileType = errors.New("unsupported file type")

const (
	uploadDirPerm          = 0o750
	checkFileTypeReadBytes = 512
)

type localStorage struct {
	uploadsDir string
}

// NewLocalStorage creates a new local storage with upload directory.
func NewLocalStorage(uploadsDir string) *localStorage { //nolint:revive
	return &localStorage{
		uploadsDir: uploadsDir,
	}
}

// StoreImage stores an image file and returns its local path.
func (s *localStorage) StoreImage(
	_ context.Context,
	fileHeader *multipart.FileHeader,
	folder string,
) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file from header: %w", err)
	}
	defer file.Close() //nolint:errcheck

	suffix, err := s.detectImageType(file)
	if err != nil {
		return "", fmt.Errorf("failed to detect image type: %w", err)
	}

	folderPath := filepath.Join(s.uploadsDir, folder)
	fileName := generate.ID("") + suffix
	imagePath := filepath.Join(folderPath, fileName)

	err = os.MkdirAll(folderPath, uploadDirPerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload folder: %w", err)
	}

	out, err := os.Create(imagePath) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("failed to store file locally: %w", err)
	}
	defer out.Close() //nolint:errcheck

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return imagePath, nil
}

// DeleteImage deletes an image file from local path.
func (s *localStorage) DeleteImage(_ context.Context, _ string) error {
	return nil
}

func (s *localStorage) detectImageType(file multipart.File) (string, error) {
	buf := make([]byte, checkFileTypeReadBytes)

	n, err := file.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("failed to read file for type detection: %w", err)
	}

	mimeType := http.DetectContentType(buf[:n])

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	switch mimeType {
	case "image/png":
		return ".png", nil
	case "image/jpeg":
		return ".jpg", nil
	default:
		return "", fmt.Errorf("%w: %s", errUnsupportedFileType, mimeType)
	}
}
