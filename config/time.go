package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func FormatDate(dateStr string) (time.Time, error) {
	t, err := time.Parse(time.DateTime, dateStr)
	return t, err
}

func FormatDateTime(dateTimeStr string) (string, error) {
	/* 	if dateTimeStr == "" {
		log.Debug().Msg("input date-time string is empty")
		return "", fmt.Errorf("input date-time string is empty")
	} */

	t, err := time.Parse(time.RFC3339, dateTimeStr)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", dateTimeStr)
		if err != nil {
			log.Err(err).Str("error parsing date-time '%s'", dateTimeStr).Msg("Error Parsing date time")

			return "", fmt.Errorf("error parsing date-time '%s': %v", dateTimeStr, err)
		}
	}

	return t.Format("2006-01-02 15:04:05"), nil
}

// Sample function to format last_update
func FormatLastUpdate(lastUpdateStr string) (string, error) {
	// Parse the original time string
	t, err := time.Parse(time.RFC3339, lastUpdateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %w", err)
	}

	// Format it to the desired layout
	formattedTime := t.Format("2006-01-02 15:04:05")
	return formattedTime, nil
}
