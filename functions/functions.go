package functions

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

var Db *bun.DB

func Contains(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsStr(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func CreateTables(ctx context.Context, db *bun.DB, models []interface{}) error {
	if db == nil {
		return fmt.Errorf("db connection is nil")
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			log.Err(err).Interface("Model:", model).Msgf("Error creating table for model")
		}

	}
	return nil
}

func ByteaToBase64(data []byte) (string, error) {
	mimeType := http.DetectContentType(data)

	if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/jpg" && mimeType != "image/gif" && mimeType != "image/svg+xml" {
		return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	base64String := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64String), nil
}

func VDecodeBase64ToByteArray(base64Image string) ([]byte, error) {
	// Check if the base64 string contains a prefix (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64Image, ";base64,"); idx != -1 {
		base64Image = base64Image[idx+8:]
	}

	// Decode the base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}

	// Return the byte array (bytea)
	return imageData, nil
}

// DecodeBase64ToByteArray decodes a base64 image string to a byte array.
func DecodeBase64ToByteArray(base64Image string) ([]byte, error) {
	// Find the start of the base64 data
	if idx := strings.Index(base64Image, "base64,"); idx != -1 {
		base64Image = base64Image[idx+7:] // Skip the 'base64,' prefix
	}

	re := regexp.MustCompile(`[^A-Za-z0-9+/=]`)
	base64Image = re.ReplaceAllString(base64Image, "")

	for len(base64Image)%4 != 0 {
		base64Image += "="
	}

	// Decode the base64 string into a byte array
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}

	return imageData, nil
}

func ParseTimeData(originalTime string) (string, error) {

	// Parse the time string to time.Time object
	parsedTime, err := time.Parse(time.RFC3339, originalTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "", err
	}

	formattedTime := parsedTime.Format("2006-01-02 15:04:05")

	return formattedTime, err
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
