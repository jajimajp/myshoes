package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/whywaita/myshoes/pkg/datastore"

	"github.com/whywaita/myshoes/internal/config"
	"github.com/whywaita/myshoes/pkg/logger"

	goji "goji.io"
	"goji.io/pat"

	httplogger "github.com/gleicon/go-httplogger"
)

// Serve start webhook receiver
func Serve(ds datastore.Datastore) error {
	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/healthz"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)

		h := struct {
			Health string `json:"health"`
		}{
			Health: "ok",
		}

		json.NewEncoder(w).Encode(h)
		return
	})

	mux.HandleFunc(pat.Post("/github/events"), func(w http.ResponseWriter, r *http.Request) {
		handleGitHubEvent(w, r, ds)
	})

	// REST API for targets
	mux.HandleFunc(pat.Post("/target"), func(w http.ResponseWriter, r *http.Request) {
		handleTargetCreate(w, r, ds)
	})
	mux.HandleFunc(pat.Get("/target/:id"), func(w http.ResponseWriter, r *http.Request) {
		handleTargetRead(w, r, ds)
	})
	mux.HandleFunc(pat.Put("/target/:id"), func(w http.ResponseWriter, r *http.Request) {
		handleTargetUpdate(w, r, ds)
	})
	mux.HandleFunc(pat.Delete("/target/:id"), func(w http.ResponseWriter, r *http.Request) {
		handleTargetDelete(w, r, ds)
	})

	logger.Logf("start webhook receiver")
	if err := http.ListenAndServe(
		":"+strconv.Itoa(config.Config.Port),
		httplogger.HTTPLogger(mux)); err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}
