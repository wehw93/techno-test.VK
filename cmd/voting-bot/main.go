package main

import (
	"log"
	"os"

	"voting-bot/internal/repository"
	"voting-bot/internal/service"
	"voting-bot/internal/transport"
)

func main() {
	token := os.Getenv("MATTERMOST_BOT_TOKEN")
	url := os.Getenv("MATTERMOST_URL")
	tarantoolHost := os.Getenv("TARANTOOL_HOST")


	repo, err := repository.NewTarantoolRepository(tarantoolHost + ":3301")
	if err != nil {
		log.Fatal("Failed to connect to Tarantool: ", err)
	}


	votingService := service.NewVotingService(repo)


	mattermostTransport, err := transport.NewMattermostTransport(url, token, votingService)
	if err != nil {
		log.Fatal("Failed to create Mattermost client: ", err)
	}


	log.Println("Bot started")
	mattermostTransport.Start()
}
