// Package handler represents entity to manipulate with request and response data
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/artnikel/apodservice/internal/constants"
	"github.com/artnikel/apodservice/internal/model"
	"github.com/sirupsen/logrus"
)

// ApodService is an interface that contains methods of service for apod
type ApodService interface {
	GetAll(ctx context.Context) ([]*model.APOD, error)
	GetByDate(ctx context.Context, date time.Time) (*model.APOD, error)
}

// ApodHandler implements methods of ApodService interface
type ApodHandler struct {
	apodSvc ApodService
}

// NewApodHandler accepts ApodService interface and returns an object of *ApodHandler
func NewApodHandler(apodSvc ApodService) *ApodHandler {
	return &ApodHandler{apodSvc: apodSvc}
}

// GetAll is a method of ApodHandler that returns all of APOD`s
// @Summary Get all APODs
// @Description Get a list of all APODs
// @ID get-all-apods
// @Produce json
// @Success 200 {array} model.APOD
// @Router /list [get]
func (h *ApodHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != constants.MethodGet {
		http.Error(w, constants.InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	apods, err := h.apodSvc.GetAll(r.Context())
	if err != nil {
		logrus.Errorf("apodHandler-getAll %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	apodsJSON, err := json.Marshal(apods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(apodsJSON)
	if err != nil {
		logrus.Println("Failed writing HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetByDate is a method of ApodHandler that returns APOD by date parameter of request
// @Summary Get APOD by date
// @Description Get APOD by the specified date
// @ID get-apod-by-date
// @Produce json
// @Param date query string true "The date in YYYY-MM-DD format"
// @Success 200 {object} model.APOD
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /bydate [get]
func (h *ApodHandler) GetByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, constants.InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "Missing 'date' parameter", http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse(constants.DateLayout, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	apod, err := h.apodSvc.GetByDate(r.Context(), parsedDate)
	if err != nil {
		logrus.Errorf("apodHandler-getByDate %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	apod.Date = dateStr
	apodJSON, err := json.Marshal(apod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(apodJSON)
	if err != nil {
		logrus.Println("Failed writing HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetToday is a method of ApodHandler that returns APOD by current date
// @Summary Get APOD for today
// @Description Get APOD for the current date
// @ID get-apod-for-today
// @Produce json
// @Success 200 {object} model.APOD
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /today [get]
func (h *ApodHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != constants.MethodGet {
		http.Error(w, constants.InvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	currentDate := time.Now().UTC()
	currentDateStr := currentDate.Format(constants.DateLayout)
	parsedCurrentDate, err := time.Parse(constants.DateLayout, currentDateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apod, err := h.apodSvc.GetByDate(r.Context(), parsedCurrentDate)
	if err != nil {
		logrus.Errorf("apodHandler-getToday %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	apod.Date = currentDateStr
	apodJSON, err := json.Marshal(apod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(apodJSON)
	if err != nil {
		logrus.Println("Failed writing HTTP response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
