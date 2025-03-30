package repository

import "voting-bot/internal/models"

type VoteRepository interface {
	CreateVote(channelID string, options []string) (string, error)
	GetVote(id string) (*models.Vote, error)
	AddVote(id string, option string) error
	EndVote(id string) error
	DeleteVote(id string) error
}
