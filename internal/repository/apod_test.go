package repository

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/model"
	"github.com/stretchr/testify/require"
)

var testApod = &model.APOD{
	Copyright:      "testCopyright",
	Date:           time.Now().UTC().Format(constants.DateLayout),
	ParsedDate:     time.Now().UTC(),
	Explanation:    "testExplanation",
	MediaType:      "testMediaType",
	ServiceVersion: "testV1",
	Title:          "testTitle",
	URL:            "testUrl",
}

func TestApodGetAllEmpty(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apods, err := client.ApodGetAll(ctx)
	require.NoError(t, err)
	require.Empty(t, apods)
}

func TestApodCreateGetAll(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testDateStr := testApod.ParsedDate.Format(constants.DateLayout)
	parsedDate, err := time.Parse(constants.DateLayout, testDateStr)
	require.NoError(t, err)
	testApod.ParsedDate = parsedDate
	err = client.ApodCreate(ctx, testApod)
	require.NoError(t, err)
	apods, err := client.ApodGetAll(ctx)
	require.NoError(t, err)
	for _, apod := range apods {
		require.Equal(t, apod.Copyright, testApod.Copyright)
		require.Equal(t, apod.Date, testApod.Date)
		require.Equal(t, apod.ParsedDate, testApod.ParsedDate)
		require.Equal(t, apod.Explanation, testApod.Explanation)
		require.Equal(t, apod.MediaType, testApod.MediaType)
		require.Equal(t, apod.ServiceVersion, testApod.ServiceVersion)
		require.Equal(t, apod.Title, testApod.Title)
		require.Equal(t, apod.URL, testApod.URL)
	}
}

func TestApodCreateGetByDate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testApod.ParsedDate = testApod.ParsedDate.Add(-24 * time.Hour)
	testApod.Date = testApod.ParsedDate.Format(constants.DateLayout)
	testDateStr := testApod.ParsedDate.Format(constants.DateLayout)
	parsedDate, err := time.Parse(constants.DateLayout, testDateStr)
	require.NoError(t, err)
	testApod.ParsedDate = parsedDate
	err = client.ApodCreate(ctx, testApod)
	require.NoError(t, err)
	currentDate := time.Now().UTC().Add(-24 * time.Hour)
	currentDateStr := currentDate.Format(constants.DateLayout)
	parsedCurrentDate, err := time.Parse(constants.DateLayout, currentDateStr)
	require.NoError(t, err)
	apod, err := client.ApodGetByDate(ctx, parsedCurrentDate)
	require.NoError(t, err)
	require.Equal(t, apod.Copyright, testApod.Copyright)
	require.Equal(t, apod.Explanation, testApod.Explanation)
	require.Equal(t, apod.MediaType, testApod.MediaType)
	require.Equal(t, apod.ServiceVersion, testApod.ServiceVersion)
	require.Equal(t, apod.Title, testApod.Title)
	require.Equal(t, apod.URL, testApod.URL)
}

func TestApodCreateWithSameDate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testApod.ParsedDate = testApod.ParsedDate.Add(-48 * time.Hour)
	testApod.Date = testApod.ParsedDate.Format(constants.DateLayout)
	testDateStr := testApod.ParsedDate.Format(constants.DateLayout)
	parsedDate, err := time.Parse(constants.DateLayout, testDateStr)
	require.NoError(t, err)
	testApod.ParsedDate = parsedDate
	err = client.ApodCreate(ctx, testApod)
	require.NoError(t, err)
	err = client.ApodCreate(ctx, testApod)
	require.Error(t, err)
}
