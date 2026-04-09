package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"beach/internal/application"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.CameraService
}

func NewHandler(service *application.CameraService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListCameras(w http.ResponseWriter, r *http.Request) {
	cameras, err := h.service.ListCameras()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorBody("internal_error"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": cameras})
}

func (h *Handler) GetCamera(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	camera, err := h.service.GetCameraByID(id)
	if err != nil {
		if errors.Is(err, application.ErrCameraNotFound) {
			writeJSON(w, http.StatusNotFound, errorBody("camera_not_found"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorBody("internal_error"))
		return
	}

	writeJSON(w, http.StatusOK, camera)
}

func (h *Handler) GetStream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	stream, err := h.service.ResolveCameraStream(id)
	if err != nil {
		if errors.Is(err, application.ErrCameraNotFound) {
			writeJSON(w, http.StatusNotFound, errorBody("camera_not_found"))
			return
		}
		if errors.Is(err, application.ErrStreamUnavailable) {
			writeJSON(w, http.StatusServiceUnavailable, errorBody("stream_unavailable"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorBody("internal_error"))
		return
	}

	writeJSON(w, http.StatusOK, stream)
}

func (h *Handler) GetConditions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cond, err := h.service.ResolveCameraConditions(id)
	if err != nil {
		if errors.Is(err, application.ErrCameraNotFound) {
			writeJSON(w, http.StatusNotFound, errorBody("camera_not_found"))
			return
		}
		if errors.Is(err, application.ErrConditionsUnavailable) {
			writeJSON(w, http.StatusServiceUnavailable, errorBody("conditions_unavailable"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorBody("internal_error"))
		return
	}

	writeJSON(w, http.StatusOK, cond)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorBody(msg string) map[string]string {
	return map[string]string{"error": msg}
}
