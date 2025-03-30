package tarantool

import (
	"strconv"

	"voting-bot/internal/models"

	"github.com/tarantool/go-tarantool"
)



type TarantoolRepository struct {
	conn *tarantool.Connection
}


func New(address string) (*TarantoolRepository, error) {
	conn, err := tarantool.Connect(address, tarantool.Opts{})
	if err != nil {
		return nil, err
	}

	return &TarantoolRepository{conn: conn}, nil
}


func (r *TarantoolRepository) CreateVote(channelID string, options []string) (string, error) {
	resp, err := r.conn.Insert("votes", []interface{}{nil, channelID, options, map[string]int{}, true})
	if err != nil {
		return "", err
	}

	id := strconv.Itoa(int(resp.Data[0].([]interface{})[0].(uint64)))
	return id, nil
}


func (r *TarantoolRepository) GetVote(id string) (*models.Vote, error) {
	resp, err := r.conn.Select("votes", "primary", 0, 1, tarantool.IterEq, []interface{}{id})
	if err != nil || len(resp.Data) == 0 {
		return nil, err
	}

	data := resp.Data[0].([]interface{})

	
	vote := &models.Vote{
		ID:        id,
		ChannelID: data[1].(string),
		Options:   make([]string, 0),
		Results:   make(map[string]int),
		Active:    data[4].(bool),
	}

	
	for _, opt := range data[2].([]interface{}) {
		vote.Options = append(vote.Options, opt.(string))
	}

	
	votesData := data[3].(map[interface{}]interface{})
	for k, v := range votesData {
		vote.Results[k.(string)] = int(v.(uint64))
	}

	return vote, nil
}


func (r *TarantoolRepository) AddVote(id string, option string) error {
	_, err := r.conn.Update("votes", "primary", []interface{}{id}, []interface{}{
		[]interface{}{"=", 4, []interface{}{[]interface{}{option, 1}}},
	})
	return err
}


func (r *TarantoolRepository) EndVote(id string) error {
	_, err := r.conn.Update("votes", "primary", []interface{}{id}, []interface{}{
		[]interface{}{"=", 4, false},
	})
	return err
}


func (r *TarantoolRepository) DeleteVote(id string) error {
	_, err := r.conn.Delete("votes", "primary", []interface{}{id})
	return err
}
