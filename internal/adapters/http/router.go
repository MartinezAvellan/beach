package http

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/index.html
var indexHTML []byte

const placeholderSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64">
<rect width="64" height="64" rx="8" fill="#1c1c2e"/>
<circle cx="32" cy="28" r="8" fill="#3b82f6" opacity="0.6"/>
<path d="M16 42 Q24 34 32 42 Q40 50 48 42" stroke="#60a5fa" stroke-width="2" fill="none" opacity="0.5"/>
</svg>`

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Frontend
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	// Fallback for missing images/svg — returns a placeholder SVG
	r.Get("/images/*", servePlaceholder)
	r.Get("/media/*", servePlaceholder)

	// API
	r.Get("/health", h.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Get("/cameras", h.ListCameras)
		r.Get("/cameras/{id}", h.GetCamera)
		r.Get("/cameras/{id}/stream", h.GetStream)
		r.Get("/cameras/{id}/conditions", h.GetConditions)
	})

	// Catch-all for any other missing static assets
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".svg") ||
			strings.HasSuffix(r.URL.Path, ".png") ||
			strings.HasSuffix(r.URL.Path, ".jpg") ||
			strings.HasSuffix(r.URL.Path, ".ico") {
			servePlaceholder(w, r)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}

func servePlaceholder(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write([]byte(placeholderSVG))
}
