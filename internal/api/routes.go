package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"bot/pkg/logger"
)

func SetupRoutes(service VotingService, logger logger.Logger) http.Handler {
	router := mux.NewRouter()
	handler := NewHandler(service, logger)


	router.Use(RecoveryMiddleware(logger))
	router.Use(LoggingMiddleware(logger))
	router.Use(ContentTypeMiddleware)
	router.Use(CORSMiddleware)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()


	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)


	pollsRouter := apiRouter.PathPrefix("/polls").Subrouter()

	pollsRouter.HandleFunc("", handler.CreatePollHandler).Methods(http.MethodPost)

	pollsRouter.HandleFunc("/{id}", handler.GetPollHandler).Methods(http.MethodGet)


	pollsRouter.HandleFunc("/{id}/results", handler.GetPollResultsHandler).Methods(http.MethodGet)


	pollsRouter.HandleFunc("/{id}/vote", handler.VoteHandler).Methods(http.MethodPost)

	pollsRouter.HandleFunc("/{id}/end", handler.EndPollHandler).Methods(http.MethodPut)


	pollsRouter.HandleFunc("/{id}", handler.DeletePollHandler).Methods(http.MethodDelete)

	return router
}