package api

import (
	"net/http"

	"github.com/gorilla/mux"

	db "github.com/tiktok/db/sqlc"
	util "github.com/tiktok/utils"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	store  db.Querier
	//tokenMaker token.Maker
	router *mux.Router
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Querier, config util.Config) (*Server, error) {

	server := &Server{
		store:  store,
		config: config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/upload/video", server.uploadVideoToS3).Methods("POST")
	router.HandleFunc("/videos", server.listAllVideos).Methods("GET")
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	err := http.ListenAndServe(address, server.router)
	if err != nil {
		return err
	}
	return nil
}
