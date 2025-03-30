package models

import (
	"time"
)

// Vote представляет голос пользователя
type Vote struct {
	UserID    string    `json:"user_id"`
	Option    string    `json:"option"`
	CreatedAt time.Time `json:"created_at"`
}

// Poll представляет голосование
type Poll struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Options     []string          `json:"options"`
	Votes       map[string][]Vote `json:"votes"` // ключ - вариант ответа, значение - список голосов
	CreatorID   string            `json:"creator_id"`
	Active      bool              `json:"active"`
	CreatedAt   time.Time         `json:"created_at"`
	EndedAt     *time.Time        `json:"ended_at,omitempty"`
}

// PollResult представляет результаты голосования
type PollResult struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Options     []string       `json:"options"`
	Results     map[string]int `json:"results"` // ключ - вариант ответа, значение - количество голосов
	TotalVotes  int            `json:"total_votes"`
	Active      bool           `json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	EndedAt     *time.Time     `json:"ended_at,omitempty"`
}

// CreatePollRequest представляет запрос на создание голосования
type CreatePollRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
}

// VoteRequest представляет запрос на голосование
type VoteRequest struct {
	UserID string `json:"user_id"`
	Option string `json:"option"`
}

// APIResponse представляет стандартный ответ API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}