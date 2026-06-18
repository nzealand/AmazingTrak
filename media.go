package main

import (
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

const maxUploadBytes = 20 << 20 // 20 MB

type UploadResult struct {
	LocalPath        string
	StoredFilename   string
	OriginalFilename string
	Lat              *float64
	Lon              *float64
	LocationSource   string
}

// processUpload reads a multipart file, re-encodes it (stripping EXIF), extracts GPS if present,
// and saves it to uploadsDir/images/. Returns metadata for the DB record.
func processUpload(fh *multipart.FileHeader, entitySlug, uploadsDir string, db *sql.DB) (*UploadResult, error) {
	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	// Read first 512 bytes to detect content type
	sniff := make([]byte, 512)
	n, _ := src.Read(sniff)
	ct := http.DetectContentType(sniff[:n])
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek: %w", err)
	}

	if !strings.HasPrefix(ct, "image/") {
		return nil, fmt.Errorf("uploaded file is not an image (detected: %s)", ct)
	}

	// Determine extension from content type
	ext := extFromContentType(ct)
	if ext == "" {
		return nil, fmt.Errorf("unsupported image type: %s", ct)
	}

	// Attempt EXIF GPS extraction before re-encoding strips it
	var lat, lon *float64
	locSrc := "unknown"
	if ct == "image/jpeg" {
		if _, err := src.Seek(0, io.SeekStart); err == nil {
			if la, lo, ok := extractGPS(src); ok {
				lat = &la
				lon = &lo
				locSrc = "exif"
			}
			src.Seek(0, io.SeekStart)
		}
	}

	// Decode image (strips all metadata including EXIF)
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek: %w", err)
	}
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	// Always store as .jpg since we re-encode to JPEG regardless of source format
	storedFilename := nextFilename(db, entitySlug, ".jpg")
	localPath := filepath.Join("images", storedFilename)
	fullPath := filepath.Join(uploadsDir, localPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	// Re-encode as JPEG (strips EXIF, normalizes format)
	if err := jpeg.Encode(dst, img, &jpeg.Options{Quality: 88}); err != nil {
		os.Remove(fullPath)
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}

	return &UploadResult{
		LocalPath:        localPath,
		StoredFilename:   storedFilename,
		OriginalFilename: fh.Filename,
		Lat:              lat,
		Lon:              lon,
		LocationSource:   locSrc,
	}, nil
}

// extractGPS reads EXIF GPS coordinates from a JPEG stream using github.com/rwcarlsen/goexif/exif.
func extractGPS(r io.Reader) (lat, lon float64, ok bool) {
	x, err := exif.Decode(r)
	if err != nil {
		return 0, 0, false
	}
	la, lo, err := x.LatLong()
	if err != nil {
		return 0, 0, false
	}
	return la, lo, true
}

// nextFilename generates a unique filename: e.g. "Amtrak-742-3.jpg"
// It counts existing stored_filename entries matching the slug prefix.
func nextFilename(db *sql.DB, entitySlug, ext string) string {
	baseName := slugToTitle(entitySlug)
	var count int
	if db != nil {
		db.QueryRow(`SELECT COUNT(*) FROM media WHERE stored_filename LIKE ?`, baseName+"-%").Scan(&count)
	}
	return fmt.Sprintf("%s-%d%s", baseName, count+1, ext)
}

// slugToTitle converts "amtrak-742" → "Amtrak-742"
func slugToTitle(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "-")
}

func extFromContentType(ct string) string {
	switch {
	case ct == "image/jpeg":
		return ".jpg"
	case ct == "image/png":
		return ".png"
	case ct == "image/gif":
		return ".gif"
	case ct == "image/webp":
		return ".webp"
	default:
		return ""
	}
}

// nextMediaSeq returns count of existing media for an entity (used to determine next seq).
func nextMediaSeq(db *sql.DB, columnName string, entityID int64) int {
	var count int
	db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM media WHERE %s=?`, columnName), entityID).Scan(&count)
	return count + 1
}

// deleteMediaFile removes a local uploaded file if it exists.
func deleteMediaFile(uploadsDir, localPath string) {
	if localPath == "" {
		return
	}
	fullPath := filepath.Join(uploadsDir, localPath)
	os.Remove(fullPath)
}

// parseOptionalLatLon parses lat/lon from form strings.
func parseOptionalLatLon(latStr, lonStr string) (lat, lon *float64) {
	latStr = strings.TrimSpace(latStr)
	lonStr = strings.TrimSpace(lonStr)
	if latStr != "" && lonStr != "" {
		var la, lo float64
		if _, err := fmt.Sscanf(latStr, "%f", &la); err == nil {
			if _, err := fmt.Sscanf(lonStr, "%f", &lo); err == nil {
				if la >= -90 && la <= 90 && lo >= -180 && lo <= 180 {
					lat = &la
					lon = &lo
				}
			}
		}
	}
	return lat, lon
}

// toNullFloat64 converts *float64 to sql.NullFloat64
func toNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

// processFaviconUpload reads an uploaded image, re-encodes it as PNG (preserving
// transparency, unlike JPEG), and saves it as a fixed filename so the site only
// ever has one favicon. Returns the path relative to uploadsDir.
func processFaviconUpload(fh *multipart.FileHeader, uploadsDir string) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	sniff := make([]byte, 512)
	n, _ := src.Read(sniff)
	ct := http.DetectContentType(sniff[:n])
	if !strings.HasPrefix(ct, "image/") {
		return "", fmt.Errorf("uploaded file is not an image (detected: %s)", ct)
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek: %w", err)
	}

	img, _, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("decode image: %w", err)
	}

	const relPath = "favicon.png"
	fullPath := filepath.Join(uploadsDir, relPath)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()
	if err := png.Encode(dst, img); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("encode png: %w", err)
	}
	return relPath, nil
}
