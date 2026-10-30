package httpapi

import (
	"context"
	"net/http"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/Zoomish/portfolio-backend/internal/content/serve"
)

type linkedinAPI interface {
	LinkedInAll(ctx context.Context) (serve.AllBody, error)
	LinkedInPart(ctx context.Context, source content.Source) (serve.PartBody, error)
}

func (h *handlers) linkedinAll(w http.ResponseWriter, r *http.Request) {
	body, err := h.linkedin.LinkedInAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, body)
}

func (h *handlers) linkedinPosts(w http.ResponseWriter, r *http.Request) {
	h.linkedinPart(w, r, content.SourceLinkedInPosts)
}

func (h *handlers) linkedinProfile(w http.ResponseWriter, r *http.Request) {
	h.linkedinPart(w, r, content.SourceLinkedInProfile)
}

func (h *handlers) linkedinContact(w http.ResponseWriter, r *http.Request) {
	h.linkedinPart(w, r, content.SourceLinkedInContact)
}

func (h *handlers) linkedinPart(w http.ResponseWriter, r *http.Request, source content.Source) {
	body, err := h.linkedin.LinkedInPart(r.Context(), source)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, body)
}
