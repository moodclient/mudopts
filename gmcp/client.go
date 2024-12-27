package gmcp

import "github.com/moodclient/telnet"

func NewPackageClient() Package {
	return Package{
		ID:      "Client",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ClientMapMessage],
			},
		},
	}
}

type ClientMapMessage struct {
	BaseMessage

	URL string `json:"url"`
}

func (m ClientMapMessage) ID() string {
	return "Client.Map"
}
