package files

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmpsec/stanza-c2/pkg/types"
	"gorm.io/gorm"
)

// ExtractedFile to keep the list of all extracted files from agents
type ExtractedFile struct {
	gorm.Model
	CommandID uint
	UUID      string `gorm:"index"`
	Fullname  string
	Size      int64
	ExfilSize int64
	MD5       string
	B64Data   string
	Verified  bool
	Extracted bool
	LocalPath string
}

// FileManager to handle all extracted files
type FileManager struct {
	DB *gorm.DB
}

// CreateFileManager to initialize the extracted file struct and its table
func CreateFileManager(backend *gorm.DB) *FileManager {
	var f *FileManager
	f = &FileManager{DB: backend}
	if err := backend.AutoMigrate(ExtractedFile{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (commands): %v", err)
	}
	return f
}

// New to convert a StzFileRequest to an ExtractedFile
func (f *FileManager) New(r *types.StzFileRequest) *ExtractedFile {
	return &ExtractedFile{
		CommandID: r.ID,
		UUID:      r.UUID,
		Fullname:  r.Fullname,
		Size:      r.Size,
		ExfilSize: r.ExfilSize,
		MD5:       r.MD5,
		B64Data:   r.B64Data,
		Verified:  false,
		Extracted: false,
		LocalPath: "",
	}
}

// Verify to check the integrity of the file
func (f *FileManager) VerifyExtract(file *ExtractedFile) ([]byte, error) {
	var data []byte
	if file.MD5 == "" {
		return data, fmt.Errorf("File %s has no MD5 hash", file.Fullname)
	}
	if file.B64Data == "" {
		return data, fmt.Errorf("File %s has no base64 data", file.Fullname)
	}
	if file.Size <= 0 {
		return data, fmt.Errorf("File %s has invalid size %d", file.Fullname, file.Size)
	}
	if file.UUID == "" {
		return data, fmt.Errorf("File %s has no UUID", file.Fullname)
	}
	if strings.HasPrefix(file.B64Data, CompressedFilePrefix) {
		encodedData := file.B64Data[5:]
		// Decode base64
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			return data, err
		}
		if !VerifyIntegrity(file.MD5, decodedData) {
			return data, fmt.Errorf("File %s integrity check failed", file.Fullname)
		}
		// Decompress with gzip
		gzipReader, err := gzip.NewReader(bytes.NewReader(decodedData))
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()
		// Read the decompressed content
		data, err = io.ReadAll(gzipReader)
		if err != nil {
			return nil, err
		}
	} else {
		if !VerifyIntegrity(file.MD5, []byte(file.B64Data)) {
			return data, fmt.Errorf("File %s integrity check failed (no compression)", file.Fullname)
		}
	}
	// If we reach this point, the file is valid and we can return the data
	if len(data) == 0 {
		return data, fmt.Errorf("File %s has no data after verification", file.Fullname)
	}
	// Last check to ensure the size matches
	if int64(len(data)) != file.Size {
		return data, fmt.Errorf("File %s size mismatch: expected %d, got %d", file.Fullname, file.Size, len(data))
	}
	log.Printf("File %s verified successfully", file.Fullname)
	return data, nil
}

// Create to create a new file in the DB
func (f *FileManager) Create(file *ExtractedFile) error {
	if err := f.DB.Create(&file).Error; err != nil {
		return fmt.Errorf("New file: %v", err)
	}
	return nil
}

// Get to retrieve a file from DB, by ID
func (f *FileManager) Get(id uint) (ExtractedFile, error) {
	var file ExtractedFile
	if err := f.DB.First(&file, id).Error; err != nil {
		return file, fmt.Errorf("Get file by ID %d: %v", id, err)
	}
	return file, nil
}

// GetByCommandID to retrieve a file from DB, by CommandID
func (f *FileManager) GetByCommandID(commandID uint) (ExtractedFile, error) {
	var file ExtractedFile
	if err := f.DB.Where("command_id = ?", commandID).First(&file).Error; err != nil {
		return file, fmt.Errorf("Get file by CommandID %d: %v", commandID, err)
	}
	return file, nil
}

// GetAll to retrieve all extracted files from DB, by UUID
func (f *FileManager) GetAll(uuid string) ([]ExtractedFile, error) {
	var files []ExtractedFile
	if err := f.DB.Find(&files).Error; err != nil {
		return files, fmt.Errorf("GetAll files by UUID %s: %v", uuid, err)
	}
	return files, nil
}

// SaveToDisk to save the file to disk
func (f *FileManager) SaveToDisk(file *ExtractedFile, data []byte) error {
	if file.LocalPath == "" {
		return fmt.Errorf("File %s has no local path to save", file.Fullname)
	}
	// Write the data to the local path
	if err := os.WriteFile(file.LocalPath, data, 0644); err != nil {
		return fmt.Errorf("Save file %s to disk: %v", file.Fullname, err)
	}
	return nil
}
