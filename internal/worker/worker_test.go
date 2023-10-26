package worker

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/artnikel/apodservice/internal/config"
	"github.com/artnikel/apodservice/internal/constants"
	"github.com/stretchr/testify/require"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetApodByKey(t *testing.T) {
	err := os.RemoveAll(constants.MediaPath)
	require.NoError(t, err)
	err = os.Mkdir(constants.MediaPath, os.FileMode(0o755))
	require.NoError(t, err)
	defer func() {
		errClose := os.RemoveAll(constants.MediaPath)
		require.NoError(t, errClose)
	}()
	apod, err := GetApodByKey(cfg)
	require.NoError(t, err)
	currentDateStr := time.Now().UTC().Format(constants.DateLayout)
	require.Equal(t, apod.Date, currentDateStr)
	require.NotEmpty(t, apod.URL)
	require.NotEmpty(t, apod.MediaType)
	require.NotEmpty(t, apod.ServiceVersion)
	require.NotEmpty(t, apod.Explanation)
	require.NotEmpty(t, apod.Title)
}

func TestDownloadMedia(t *testing.T) {
	var (
		testMediaType = "image"
		testAPIURL    = cfg.NasaAPIURL + cfg.NasaAPIKey
		testDate      = time.Now().UTC().Format(constants.DateLayout)
	)
	err := os.RemoveAll(constants.MediaPath)
	require.NoError(t, err)
	err = os.Mkdir(constants.MediaPath, os.FileMode(0o755))
	require.NoError(t, err)
	err = downloadMedia(testMediaType, testAPIURL, testDate)
	require.NoError(t, err)
	_, err = os.Stat(constants.MediaPath + "/" + testDate + ".jpg")
	require.NoError(t, err)
	err = os.RemoveAll(constants.MediaPath)
	require.NoError(t, err)
	testMediaType = "notImage"
	err = downloadMedia(testMediaType, testAPIURL, testDate)
	require.NoError(t, err)
	_, err = os.Stat(constants.MediaPath + "/" + testDate + ".jpg")
	require.Error(t, err)
}
