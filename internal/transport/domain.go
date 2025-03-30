package transport

import (
	"voting-bot/internal/service"

	"github.com/mattermost/mattermost-server/v6/model"
)
type MattermostTransport struct {
	client          *model.Client4
	user            *model.User
	webSocket       *model.WebSocketClient
	votingService   service.VotingService
	commandHandlers map[string]func(string, string)
}
func (t *MattermostTransport) Start() {
	t.webSocket.Listen()

	for event := range t.webSocket.EventChannel {
		t.handleEvent(event)
	}
}

func NewMattermostTransport(url, token string, votingService service.VotingService) (*MattermostTransport, error) {
	client := model.NewAPIv4Client(url)
	client.SetToken(token)

	user, _, err := client.GetMe("")
	if err != nil {
		return nil, err
	}

	ws, err := model.NewWebSocketClient4(url, token)
	if err != nil {
		return nil, err
	}

	transport := &MattermostTransport{
		client:        client,
		user:          user,
		webSocket:     ws,
		votingService: votingService,
	}

	transport.commandHandlers = map[string]func(string, string){
		"create":  transport.handleCreate,
		"vote":    transport.handleVote,
		"results": transport.handleResults,
		"end":     transport.handleEnd,
		"delete":  transport.handleDelete,
	}

	return transport, nil
}
