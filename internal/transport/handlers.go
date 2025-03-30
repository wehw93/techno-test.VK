package transport

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)



func (t *MattermostTransport) handleEvent(event *model.WebSocketEvent) {
	if event.EventType() != model.WebsocketEventPosted {
		return
	}

	var post model.Post
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil || post.UserId == t.user.Id {
		return
	}

	message := strings.TrimSpace(post.Message)
	if !strings.HasPrefix(message, "/vote") {
		return
	}

	parts := strings.Fields(message)
	if len(parts) < 2 {
		return
	}

	command := parts[1]
	if handler, ok := t.commandHandlers[command]; ok {
		handler(post.ChannelId, strings.Join(parts[2:], " "))
	}
}

func (t *MattermostTransport) handleCreate(channelID, args string) {
	id, options, err := t.votingService.CreateVoting(channelID, args)
	if err != nil {
		t.sendMessage(channelID, "Ошибка: "+err.Error()+"\nИспользование: /vote create option1|option2|option3")
		return
	}

	msg := "Голосование создано с ID: " + id + "\nВарианты: " + strings.Join(options, ", ")
	t.sendMessage(channelID, msg)
}

func (t *MattermostTransport) handleVote(channelID, args string) {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) != 2 {
		t.sendMessage(channelID, "Использование: /vote vote <id> <option>")
		return
	}

	id := parts[0]
	option := parts[1]

	err := t.votingService.RecordVote(id, option)
	if err != nil {
		t.sendMessage(channelID, "Ошибка: "+err.Error())
		return
	}

	t.sendMessage(channelID, "Голос записан")
}

func (t *MattermostTransport) handleResults(channelID, id string) {
	vote, err := t.votingService.GetResults(id)
	if err != nil {
		t.sendMessage(channelID, "Голосование не найдено")
		return
	}

	var result strings.Builder
	result.WriteString("Результаты голосования (ID: " + id + "):\n")
	for _, option := range vote.Options {
		count := vote.Results[option]
		result.WriteString(option + ": " + strconv.Itoa(count) + "\n")
	}

	if vote.Active {
		result.WriteString("Статус: Активно")
	} else {
		result.WriteString("Статус: Закрыто")
	}

	t.sendMessage(channelID, result.String())
}

func (t *MattermostTransport) handleEnd(channelID, id string) {
	err := t.votingService.EndVoting(id)
	if err != nil {
		t.sendMessage(channelID, "Ошибка при завершении голосования")
		return
	}

	t.sendMessage(channelID, "Голосование завершено")
}

func (t *MattermostTransport) handleDelete(channelID, id string) {
	err := t.votingService.DeleteVoting(id)
	if err != nil {
		t.sendMessage(channelID, "Ошибка при удалении голосования")
		return
	}

	t.sendMessage(channelID, "Голосование удалено")
}

func (t *MattermostTransport) sendMessage(channelID, message string) {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}

	_, _, err := t.client.CreatePost(post)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
