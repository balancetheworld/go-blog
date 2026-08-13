package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

const MaxImageUploadSize int64 = 10 * 1024 * 1024

var ErrImageTooLarge = errors.New("image exceeds size limit")
var ErrInvalidImage = errors.New("invalid image file")

var imageExtensions = map[string]map[string]bool{
	"image/gif":  {".gif": true},
	"image/jpeg": {".jpeg": true, ".jpg": true},
	"image/png":  {".png": true},
	"image/webp": {".webp": true},
}

type UploadedImage struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func FileBasePath() string {
	return Get(constant.EnvKeyFileBasePath, "data/uploads")
}

func ResolveImagePath(value string) (string, error) {
	basePath, err := filepath.Abs(FileBasePath())
	if err != nil {
		return "", err
	}

	relativePath := strings.TrimPrefix(filepath.Clean("/"+value), "/")
	if relativePath == "" || relativePath == "." {
		return "", ErrInvalidImage
	}

	resolvedPath := filepath.Join(basePath, relativePath)
	if !strings.HasPrefix(resolvedPath, basePath+string(os.PathSeparator)) {
		return "", ErrInvalidImage
	}

	return resolvedPath, nil
}

func SaveImage(fileHeader *multipart.FileHeader) (UploadedImage, error) {
	if fileHeader == nil {
		return UploadedImage{}, ErrInvalidImage
	}
	if fileHeader.Size <= 0 {
		return UploadedImage{}, ErrInvalidImage
	}
	if fileHeader.Size > MaxImageUploadSize {
		return UploadedImage{}, ErrImageTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return UploadedImage{}, err
	}
	defer file.Close()

	header := make([]byte, 512)
	readSize, err := io.ReadFull(file, header)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return UploadedImage{}, ErrInvalidImage
	}

	contentType := http.DetectContentType(header[:readSize])
	extensions, ok := imageExtensions[contentType]
	if !ok {
		return UploadedImage{}, ErrInvalidImage
	}

	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !extensions[extension] {
		return UploadedImage{}, ErrInvalidImage
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return UploadedImage{}, err
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return UploadedImage{}, err
	}

	now := time.Now()
	relativeDirectory := filepath.Join(
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
	)
	filename := hex.EncodeToString(randomBytes) + extension
	directory := filepath.Join(FileBasePath(), relativeDirectory)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return UploadedImage{}, err
	}

	destination, err := os.OpenFile(
		filepath.Join(directory, filename),
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		return UploadedImage{}, err
	}

	written, copyErr := io.Copy(
		destination,
		io.LimitReader(file, MaxImageUploadSize+1),
	)
	closeErr := destination.Close()
	if copyErr != nil {
		return UploadedImage{}, copyErr
	}
	if closeErr != nil {
		return UploadedImage{}, closeErr
	}
	if written > MaxImageUploadSize {
		_ = os.Remove(filepath.Join(directory, filename))
		return UploadedImage{}, ErrImageTooLarge
	}

	urlPath := strings.Join([]string{
		"/api/v1/file/content",
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		filename,
	}, "/")

	return UploadedImage{
		URL:  urlPath,
		Name: filename,
		Size: written,
	}, nil
}
