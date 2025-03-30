package service

import (
	"errors"
	"strings"

	"voting-bot/internal/repository"
	"voting-bot/internal/models"
)


type VotingService interface {
	CreateVoting(channelID string, args string) (string, []string, error)
	RecordVote(id string, option string) error
	GetResults(id string) (*models.Vote, error)
	EndVoting(id string) error
	DeleteVoting(id string) error
	IsActive(id string) bool
}


type VotingServiceImpl struct {
	repo          repository.VoteRepository
	activeVotings map[string]bool
}


func NewVotingService(repo repository.VoteRepository) *VotingServiceImpl {
	return &VotingServiceImpl{
		repo:          repo,
		activeVotings: make(map[string]bool),
	}
}


func (s *VotingServiceImpl) CreateVoting(channelID string, args string) (string, []string, error) {
	options := strings.Split(args, "|")
	if len(options) < 2 {
		return "", nil, errors.New("at least two options required")
	}

	for i := range options {
		options[i] = strings.TrimSpace(options[i])
	}

	id, err := s.repo.CreateVote(channelID, options)
	if err != nil {
		return "", nil, err
	}

	s.activeVotings[id] = true
	return id, options, nil
}

func (s *VotingServiceImpl) RecordVote(id string, option string) error {
	if !s.IsActive(id) {
		return errors.New("voting is not active")
	}

	return s.repo.AddVote(id, option)
}


func (s *VotingServiceImpl) GetResults(id string) (*models.Vote, error) {
	return s.repo.GetVote(id)
}


func (s *VotingServiceImpl) EndVoting(id string) error {
	err := s.repo.EndVote(id)
	if err != nil {
		return err
	}

	delete(s.activeVotings, id)
	return nil
}


func (s *VotingServiceImpl) DeleteVoting(id string) error {
	err := s.repo.DeleteVote(id)
	if err != nil {
		return err
	}

	delete(s.activeVotings, id)
	return nil
}

func (s *VotingServiceImpl) IsActive(id string) bool {
	active, exists := s.activeVotings[id]
	return exists && active
}
