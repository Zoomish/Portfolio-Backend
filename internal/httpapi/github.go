package httpapi

import (
	"context"
	"net/http"

	"github.com/Zoomish/portfolio-backend/internal/content/serve"
)

type githubAPI interface {
	Github(ctx context.Context) (serve.GithubBody, error)
}

func (h *handlers) githubRepos(w http.ResponseWriter, r *http.Request) {
	body, err := h.github.Github(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, body)
}
