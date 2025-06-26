package files

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// CompressedFilePrefix is the prefix for compressed files
	CompressedFilePrefix = "GZIP:"
)

// Helper to get the MD5 hash of a blob
func GetMD5(blob []byte) string {
	if len(blob) == 0 {
		return ""
	}
	h := md5.New()
	h.Write(blob)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Helper to verify the integrity of a file using its MD5 hash
func VerifyIntegrity(fileMD5 string, fileData []byte) bool {
	if fileMD5 == "" {
		return false
	}
	if len(fileData) == 0 {
		return false
	}
	return (GetMD5(fileData) == fileMD5)
}

// Helper to generate a file name from a UUID
func ProcessedFilename(uuid, filename string) string {
	now := time.Now()
	timestamp := now.Unix()
	return fmt.Sprintf("%s-%d-%s", uuid, timestamp, ProcessedFilenameFromFile(filename))
}

// Helper to convert a filename to a processed filename
func ProcessedFilenameFromFile(filename string) string {
	if filename == "" {
		return ""
	}
	// Replace spaces and slashes with underscores
	processed := filename
	// Replace spaces with underscores
	processed = strings.ReplaceAll(processed, " ", "_")
	// Replace forward and backward slashes with underscores
	processed = strings.ReplaceAll(processed, "/", "_")
	processed = strings.ReplaceAll(processed, "\\", "_")
	// Remove special characters (keep alphanumeric, underscores, periods, and hyphens)
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-\.]`)
	processed = reg.ReplaceAllString(processed, "")
	return processed
}
