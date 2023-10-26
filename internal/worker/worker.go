// Package worker contains methods for sending a request to NASA
package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/artnikel/apodservice/internal/config"
	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/model"
	"github.com/sirupsen/logrus"
)

// GetApodByKey is method of worker that get metadata from NASA  by API Key to struct
func GetApodByKey(cfg *config.Config) (*model.APOD, error) {
	apiURL := cfg.NasaAPIURL + cfg.NasaAPIKey
	resp, err := http.Get(apiURL) // nolint
	if err != nil {
		return nil, fmt.Errorf("get %v", err)
	}
	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			logrus.Errorf("close %v", errClose)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	var apod *model.APOD
	err = json.NewDecoder(resp.Body).Decode(&apod)
	if err != nil {
		return nil, fmt.Errorf("decode %v", err)
	}
	parsedDate, dateParseErr := time.Parse(constants.DateLayout, apod.Date)
	if dateParseErr != nil {
		return nil, fmt.Errorf("parse %v", dateParseErr)
	}
	apod.ParsedDate = parsedDate
	err = downloadMedia(apod.MediaType, apod.URL, apod.Date)
	if err != nil {
		return nil, fmt.Errorf("downloadMedia %w", err)
	}
	return apod, nil
}

// downloadMedia save image to storage
func downloadMedia(mediaType, url, date string) error {
	resp, err := http.Get(url) // nolint
	if err != nil {
		return fmt.Errorf("get %v", err)
	}
	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			logrus.Errorf("close %v", errClose)
		}
	}()
	if mediaType == "image" {
		{
			mediaName := date + ".jpg"
			dst, err := os.Create(filepath.Join(constants.MediaPath, mediaName)) // nolint
			if err != nil {
				return fmt.Errorf("create %v", err)
			}
			defer func() {
				errClose := dst.Close()
				if errClose != nil {
					logrus.Errorf("close %v", errClose)
				}
			}()
			_, err = io.Copy(dst, resp.Body)
			if err != nil {
				return fmt.Errorf("copy %v", err)
			}
		}
	}
	return nil
}
