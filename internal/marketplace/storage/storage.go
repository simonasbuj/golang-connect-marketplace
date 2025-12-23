// Package storage handles local storage of listing images.
package storage

import (
	"context"
	"mime/multipart"
)

// Storage defines methods for storing and deleting listing images.
type Storage interface {
	StoreImage(ctx context.Context, fileHeader *multipart.FileHeader, path string) (string, error)
	DeleteImage(ctx context.Context, path string) error
}
