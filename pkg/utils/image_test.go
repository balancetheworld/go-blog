package utils

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zyj/my-blog/pkg/constant"
)

func imageFileHeader(t *testing.T, name, contentType string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", `form-data; name="file"; filename="`+name+`"`)
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader := multipart.NewReader(&body, writer.Boundary())
	form, err := reader.ReadForm(int64(len(content) + 1024))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })
	return form.File["file"][0]
}

func TestSaveImage(t *testing.T) {
	uploadDirectory := t.TempDir()
	t.Setenv(constant.EnvKeyFileBasePath, uploadDirectory)
	content := append(
		[]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a},
		make([]byte, 512)...,
	)

	result, err := SaveImage(imageFileHeader(t, "test.png", "image/png", content))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result.URL, "/api/v1/file/content/") {
		t.Fatalf("unexpected image URL: %s", result.URL)
	}
	if result.Size != int64(len(content)) {
		t.Fatalf("unexpected image size: %d", result.Size)
	}

	storedFiles := 0
	err = filepath.Walk(uploadDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			storedFiles++
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if storedFiles != 1 {
		t.Fatalf("unexpected stored file count: %d", storedFiles)
	}
}

func TestSaveImageRejectsInvalidContent(t *testing.T) {
	t.Setenv(constant.EnvKeyFileBasePath, t.TempDir())
	_, err := SaveImage(imageFileHeader(
		t,
		"test.png",
		"image/png",
		[]byte("not an image"),
	))
	if err != ErrInvalidImage {
		t.Fatalf("expected invalid image error, got %v", err)
	}
}

func TestResolveImagePath(t *testing.T) {
	uploadDirectory := t.TempDir()
	t.Setenv(constant.EnvKeyFileBasePath, uploadDirectory)

	resolvedPath, err := ResolveImagePath("2026/08/image.png")
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := filepath.Join(uploadDirectory, "2026", "08", "image.png")
	if resolvedPath != expectedPath {
		t.Fatalf("expected %s, got %s", expectedPath, resolvedPath)
	}
}
