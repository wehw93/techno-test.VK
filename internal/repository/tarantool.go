package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"

	"bot/internal/models"
	"bot/pkg/logger"
)


type TarantoolRepository struct {
	conn   *tarantool.Connection
	logger logger.Logger
}


func NewTarantoolRepository(addr string, logger logger.Logger) (*TarantoolRepository, error) {
	opts := tarantool.Opts{
		Timeout:       5 * time.Second,
		Reconnect:     1 * time.Second,
		MaxReconnects: 10,
	}

	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to Tarantool at " + addr)


	if err := initSchema(conn, logger); err != nil {
		conn.Close()
		return nil, err
	}

	return &TarantoolRepository{
		conn:   conn,
		logger: logger,
	}, nil
}

                  
func (r *TarantoolRepository) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

                  
func initSchema(conn *tarantool.Connection, logger logger.Logger) error {

	_, err := conn.Call("box.schema.space.create", []interface{}{"polls", map[string]bool{"if_not_exists": true}})
	if err != nil {
		logger.Error("Failed to create polls space", err)
		return err
	}

	_, err = conn.Call("box.space.polls:create_index", []interface{}{
		"primary",
		map[string]interface{}{
			"type":          "HASH",
			"parts":         []interface{}{1}, 
			"if_not_exists": true,
		},
	})
	if err != nil {
		logger.Error("Failed to create primary index for polls space", err)
		return err
	}

	logger.Info("Tarantool schema initialized successfully")
	return nil
}


func (r *TarantoolRepository) CreatePoll(poll models.Poll) (string, error) {
	if poll.ID == "" {
		poll.ID = uuid.New().String()
	}

	pollData, err := packPoll(poll)
	if err != nil {
		return "", err
	}

	_, err = r.conn.Insert("polls", []interface{}{poll.ID, pollData})
	if err != nil {
		r.logger.Error("Failed to create poll", err)
		return "", err
	}

	r.logger.Info(fmt.Sprintf("Created poll with ID: %s", poll.ID))
	return poll.ID, nil
}


func (r *TarantoolRepository) GetPoll(id string) (models.Poll, error) {
	resp, err := r.conn.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{id})
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to get poll with ID: %s", id), err)
		return models.Poll{}, err
	}

	if len(resp.Data) == 0 {
		return models.Poll{}, errors.New("poll not found")
	}

	tuples := resp.Tuples()
	if len(tuples) == 0 || len(tuples[0]) < 2 {
		return models.Poll{}, errors.New("invalid poll data")
	}

	pollData, ok := tuples[0][1].(map[string]interface{})
	if !ok {
		return models.Poll{}, errors.New("invalid poll data format")
	}

	return unpackPoll(pollData)
}


func (r *TarantoolRepository) UpdatePoll(poll models.Poll) error {
	pollData, err := packPoll(poll)
	if err != nil {
		return err
	}

	_, err = r.conn.Replace("polls", []interface{}{poll.ID, pollData})
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to update poll with ID: %s", poll.ID), err)
		return err
	}

	r.logger.Info(fmt.Sprintf("Updated poll with ID: %s", poll.ID))
	return nil
}

func (r *TarantoolRepository) DeletePoll(id string) error {
	_, err := r.conn.Delete("polls", "primary", []interface{}{id})
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to delete poll with ID: %s", id), err)
		return err
	}

	r.logger.Info(fmt.Sprintf("Deleted poll with ID: %s", id))
	return nil
}


func packPoll(poll models.Poll) (map[string]interface{}, error) {

	votes := make(map[string]interface{})
	for option, optionVotes := range poll.Votes {
		votesData := make([]map[string]interface{}, len(optionVotes))
		for i, vote := range optionVotes {
			votesData[i] = map[string]interface{}{
				"user_id":    vote.UserID,
				"option":     vote.Option,
				"created_at": vote.CreatedAt.Unix(),
			}
		}
		votes[option] = votesData
	}


	pollData := map[string]interface{}{
		"id":          poll.ID,
		"title":       poll.Title,
		"description": poll.Description,
		"options":     poll.Options,
		"votes":       votes,
		"creator_id":  poll.CreatorID,
		"active":      poll.Active,
		"created_at":  poll.CreatedAt.Unix(),
	}

	if poll.EndedAt != nil {
		pollData["ended_at"] = poll.EndedAt.Unix()
	}

	return pollData, nil
}

func unpackPoll(data map[string]interface{}) (models.Poll, error) {
	poll := models.Poll{
		ID:          data["id"].(string),
		Title:       data["title"].(string),
		Description: data["description"].(string),
		CreatorID:   data["creator_id"].(string),
		Active:      data["active"].(bool),
		CreatedAt:   time.Unix(data["created_at"].(int64), 0),
		Votes:       make(map[string][]models.Vote),
	}


	optionsData, ok := data["options"].([]interface{})
	if !ok {
		return models.Poll{}, errors.New("invalid options format")
	}

	poll.Options = make([]string, len(optionsData))
	for i, opt := range optionsData {
		poll.Options[i] = opt.(string)
	}


	votesData, ok := data["votes"].(map[string]interface{})
	if ok {
		for option, votes := range votesData {
			optionVotes, ok := votes.([]interface{})
			if !ok {
				continue
			}

			poll.Votes[option] = make([]models.Vote, len(optionVotes))
			for i, voteData := range optionVotes {
				voteMap, ok := voteData.(map[string]interface{})
				if !ok {
					continue
				}

				poll.Votes[option][i] = models.Vote{
					UserID:    voteMap["user_id"].(string),
					Option:    voteMap["option"].(string),
					CreatedAt: time.Unix(voteMap["created_at"].(int64), 0),
				}
			}
		}
	}

	if endedAt, ok := data["ended_at"].(int64); ok {
		t := time.Unix(endedAt, 0)
		poll.EndedAt = &t
	}

	return poll, nil
}