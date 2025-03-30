package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"bot/internal/models"
	"bot/pkg/logger"
)


type VotingRepository interface {
	CreatePoll(poll models.Poll) (string, error)
	GetPoll(id string) (models.Poll, error)
	UpdatePoll(poll models.Poll) error
	DeletePoll(id string) error
}


type VotingService struct {
	repo   VotingRepository
	logger logger.Logger
}


func NewVotingService(repo VotingRepository, logger logger.Logger) *VotingService {
	return &VotingService{
		repo:   repo,
		logger: logger,
	}
}


func (s *VotingService) CreatePoll(req models.CreatePollRequest, creatorID string) (models.Poll, error) {

	if req.Title == "" {
		return models.Poll{}, errors.New("title is required")
	}
	if len(req.Options) < 2 {
		return models.Poll{}, errors.New("at least two options are required")
	}


	poll := models.Poll{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Options:     req.Options,
		Votes:       make(map[string][]models.Vote),
		CreatorID:   creatorID,
		Active:      true,
		CreatedAt:   time.Now(),
	}

	for _, option := range req.Options {
		poll.Votes[option] = []models.Vote{}
	}


	pollID, err := s.repo.CreatePoll(poll)
	if err != nil {
		s.logger.Error("Error creating poll", err)
		return models.Poll{}, err
	}

	poll.ID = pollID
	s.logger.Info("Poll created successfully with ID: " + pollID)
	return poll, nil
}


func (s *VotingService) GetPoll(id string) (models.Poll, error) {
	poll, err := s.repo.GetPoll(id)
	if err != nil {
		s.logger.Error("Error getting poll", err)
		return models.Poll{}, err
	}
	return poll, nil
}


func (s *VotingService) GetPollResults(id string) (models.PollResult, error) {
	poll, err := s.repo.GetPoll(id)
	if err != nil {
		s.logger.Error("Error getting poll for results", err)
		return models.PollResult{}, err
	}

	results := make(map[string]int)
	totalVotes := 0
	for option, votes := range poll.Votes {
		results[option] = len(votes)
		totalVotes += len(votes)
	}

	pollResult := models.PollResult{
		ID:          poll.ID,
		Title:       poll.Title,
		Description: poll.Description,
		Options:     poll.Options,
		Results:     results,
		TotalVotes:  totalVotes,
		Active:      poll.Active,
		CreatedAt:   poll.CreatedAt,
		EndedAt:     poll.EndedAt,
	}

	return pollResult, nil
}

func (s *VotingService) Vote(pollID string, req models.VoteRequest) error {

	poll, err := s.repo.GetPoll(pollID)
	if err != nil {
		s.logger.Error("Error getting poll for voting", err)
		return err
	}

	
	if !poll.Active {
		return errors.New("poll is not active")
	}

	
	optionValid := false
	for _, option := range poll.Options {
		if option == req.Option {
			optionValid = true
			break
		}
	}
	if !optionValid {
		return errors.New("invalid option")
	}


	for _, votes := range poll.Votes {
		for _, vote := range votes {
			if vote.UserID == req.UserID {
				return errors.New("user has already voted")
			}
		}
	}


	vote := models.Vote{
		UserID:    req.UserID,
		Option:    req.Option,
		CreatedAt: time.Now(),
	}
	poll.Votes[req.Option] = append(poll.Votes[req.Option], vote)

	if err := s.repo.UpdatePoll(poll); err != nil {
		s.logger.Error("Error updating poll after vote", err)
		return err
	}

	s.logger.Info("Vote registered for poll " + pollID + " by user " + req.UserID)
	return nil
}


func (s *VotingService) EndPoll(pollID string, userID string) error {

	poll, err := s.repo.GetPoll(pollID)
	if err != nil {
		s.logger.Error("Error getting poll for ending", err)
		return err
	}


	if poll.CreatorID != userID {
		return errors.New("only the creator can end the poll")
	}

	
	if !poll.Active {
		return errors.New("poll is already ended")
	}


	poll.Active = false
	now := time.Now()
	poll.EndedAt = &now

	
	if err := s.repo.UpdatePoll(poll); err != nil {
		s.logger.Error("Error updating poll after ending", err)
		return err
	}

	s.logger.Info("Poll " + pollID + " ended by user " + userID)
	return nil
}


func (s *VotingService) DeletePoll(pollID string, userID string) error {
	
	poll, err := s.repo.GetPoll(pollID)
	if err != nil {
		s.logger.Error("Error getting poll for deletion", err)
		return err
	}


	if poll.CreatorID != userID {
		return errors.New("only the creator can delete the poll")
	}


	if err := s.repo.DeletePoll(pollID); err != nil {
		s.logger.Error("Error deleting poll", err)
		return err
	}

	s.logger.Info("Poll " + pollID + " deleted by user " + userID)
	return nil
}