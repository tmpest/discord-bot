package server

import (
	"fmt"
	"net/http"
	"os"
)

// The DiscordBotServer that encapsulates logic for the DiscordBot
type DiscordBotServer struct {
	server *http.Server
}

// New constructs a new DiscordBotServer that is deployable on Heroku
func New() *DiscordBotServer {
	return &DiscordBotServer{
		server: &http.Server{
			//Use the port Heroku assigns via $PORT, all web traffic to the URL specified in the Heroku App Settings page
			//is redirected to our server on the specified port. Probably sharing resources
			Addr:    fmt.Sprintf(":%+v", os.Getenv("PORT")),
			Handler: oAuth2RedirectHandler{},
		},
	}
}

// Start turns the server on and keeps running until killed
func (s *DiscordBotServer) Start() {
	fmt.Printf("Starting the Discord Bot Server at: %+v\n", s.server.Addr)
	defer fmt.Println("Exiting! Goodbye!")

	s.server.ListenAndServe()
}
