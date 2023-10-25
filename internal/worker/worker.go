package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/artnikel/apodservice/internal/config"
	"github.com/artnikel/apodservice/internal/model"
)

func GetApod(cfg *config.Config) (*model.APOD, error) {
	apiURL := cfg.NasaApiUrl + cfg.NasaApiKey
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	var apod *model.APOD
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return nil, err
	}
	dateLayout := "2006-01-02"
	parsedDate, dateParseErr := time.Parse(dateLayout, apod.Date)
	if dateParseErr != nil {
		return nil, fmt.Errorf("ошибка при разборе даты: %v", dateParseErr)
	}
	apod.ParsedDate = parsedDate

	return apod, nil
}
