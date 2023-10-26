package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/handler/mocks"
	"github.com/artnikel/apodservice/internal/model"
	"github.com/stretchr/testify/mock"
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

func TestGetAll(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	testApods := make([]*model.APOD, 0)
	testApods = append(testApods, testApod)
	srv.On("GetAll", mock.Anything).Return(testApods, nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/list", http.NoBody)
	hndl.GetAll(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	apods := make([]*model.APOD, 0)
	err := json.NewDecoder(resp.Body).Decode(&apods)
	require.NoError(t, err)
	for _, apod := range apods {
		require.Equal(t, apod.Title, testApod.Title)
		require.Equal(t, apod.Copyright, testApod.Copyright)
		require.Equal(t, apod.MediaType, testApod.MediaType)
		require.Equal(t, apod.URL, testApod.URL)
		require.Equal(t, apod.Explanation, testApod.Explanation)
	}
}

func TestGetAllEmpty(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	srv.On("GetAll", mock.Anything).Return(nil, errors.New("error")).Once()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/list", http.NoBody)
	hndl.GetAll(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestGetByDate(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	srv.On("GetByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(testApod, nil).Once()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/bydate?date="+testApod.Date, http.NoBody)
	hndl.GetByDate(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	var apod *model.APOD
	err := json.NewDecoder(resp.Body).Decode(&apod)
	require.NoError(t, err)
	require.Equal(t, apod.Title, testApod.Title)
	require.Equal(t, apod.Copyright, testApod.Copyright)
	require.Equal(t, apod.MediaType, testApod.MediaType)
	require.Equal(t, apod.URL, testApod.URL)
	require.Equal(t, apod.Explanation, testApod.Explanation)
}

func TestGetByDateNotFound(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	srv.On("GetByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(nil, errors.New("error")).Once()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/bydate?date="+testApod.Date, http.NoBody)
	hndl.GetByDate(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestGetByDateWithoutParameters(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/bydate", http.NoBody)
	hndl.GetByDate(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusBadRequest)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/bydate?date=invalid_date", http.NoBody)
	hndl.GetByDate(w, r)
	resp = w.Result()
	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestGetToday(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	srv.On("GetByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(testApod, nil).Once()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/today", http.NoBody)
	hndl.GetToday(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	var apod *model.APOD
	err := json.NewDecoder(resp.Body).Decode(&apod)
	require.NoError(t, err)
	require.Equal(t, apod.Title, testApod.Title)
	require.Equal(t, apod.Copyright, testApod.Copyright)
	require.Equal(t, apod.MediaType, testApod.MediaType)
	require.Equal(t, apod.URL, testApod.URL)
	require.Equal(t, apod.Explanation, testApod.Explanation)
}

func TestGetTodayNotFound(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	srv.On("GetByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(nil, errors.New("error")).Once()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/today", http.NoBody)
	hndl.GetToday(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestMethodsWrong(t *testing.T) {
	srv := new(mocks.ApodService)
	hndl := NewApodHandler(srv)
	testApods := make([]*model.APOD, 0)
	testApods = append(testApods, testApod)
	srv.On("GetAll", mock.Anything).Return(testApods, nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/list", http.NoBody)
	hndl.GetAll(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), constants.InvalidMethod+"\n")

	srv.On("GetByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(testApod, nil).Once()
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodDelete, "/bydate?date="+testApod.Date, http.NoBody)
	hndl.GetByDate(w, r)
	resp = w.Result()
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), constants.InvalidMethod+"\n")

	srv.On("GetByDate", mock.Anything, mock.AnythingOfType("time.Time")).Return(testApod, nil).Once()
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPut, "/today", http.NoBody)
	hndl.GetByDate(w, r)
	resp = w.Result()
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), constants.InvalidMethod+"\n")
}
