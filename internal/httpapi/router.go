package httpapi

import "net/http"

type handlers struct {
	github   githubAPI
	linkedin linkedinAPI
	ready    readyCheck
}

func New(github githubAPI, linkedin linkedinAPI, ready readyCheck) http.Handler {
	h := &handlers{github: github, linkedin: linkedin, ready: ready}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/github/repos", h.githubRepos)
	mux.HandleFunc("GET /api/linkedin/all", h.linkedinAll)
	mux.HandleFunc("GET /api/linkedin/posts", h.linkedinPosts)
	mux.HandleFunc("GET /api/linkedin/profile", h.linkedinProfile)
	mux.HandleFunc("GET /api/linkedin/contact", h.linkedinContact)
	mux.HandleFunc("GET /api/app", h.app)
	mux.HandleFunc("GET /api/ready", h.readiness)
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
