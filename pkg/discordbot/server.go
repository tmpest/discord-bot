package discordbot

import (
	"fmt"
	"net/http"
	"os"
)

// Server to provide features to a Bot
type Server struct {
	server *http.Server
}

// NewServer constructs a Server configured Heroku
func NewServer() *Server {
	return &Server{
		server: &http.Server{
			//Use the port Heroku assigns via $PORT, all web traffic to the URL specified in the Heroku App Settings page
			//is redirected to our server on the specified port. Probably sharing resources
			Addr:    fmt.Sprintf(":%+v", os.Getenv("PORT")),
			Handler: oAuth2RedirectHandler{},
		},
	}
}

// Start turns the server on and keeps running until killed
func (s *Server) Start() {
	fmt.Printf("Starting the Discord Bot Server at: %+v\n", s.server.Addr)
	defer fmt.Println("Exiting! Goodbye!")

	s.server.ListenAndServe()
}
