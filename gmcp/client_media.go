package gmcp

import "github.com/moodclient/telnet"

func NewPackageClientMedia() Package {
	return Package{
		ID:      "Client.Media",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ClientMediaDefaultMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ClientMediaLoadMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ClientMediaPlayMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ClientMediaStopMessage],
			},
		},
	}
}

type ClientMediaDefaultMessage struct {
	BaseMessage

	URL string `json:"url"`
}

func (m ClientMediaDefaultMessage) ID() string {
	return "Client.Media.Default"
}

type ClientMediaLoadMessage struct {
	BaseMessage

	Name string `json:"name"`
	URL  string `json:"url"`
}

func (m ClientMediaLoadMessage) ID() string {
	return "Client.Media.Load"
}

type ClientMediaPlayMessage struct {
	BaseMessage

	Name     string `json:"name"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Tag      string `json:"tag"`
	Volume   int    `json:"volume"`
	FadeIn   int    `json:"fadein"`
	FadeOut  int    `json:"fadeout"`
	Start    int    `json:"start"`
	Finish   int    `json:"finish"`
	Loops    int    `json:"loops"`
	Priority int    `json:"priority"`
	Continue bool   `json:"continue"`
	Key      string `json:"key"`
}

func (m ClientMediaPlayMessage) ID() string {
	return "Client.Media.Play"
}

type ClientMediaStopMessage struct {
	BaseMessage

	Name     string `json:"name"`
	Type     string `json:"type"`
	Tag      string `json:"tag"`
	Priority int    `json:"priority"`
	Key      string `json:"key"`
	FadeAway bool   `json:"fadeaway"`
	FadeOut  int    `json:"fadeout"`
}

func (m ClientMediaStopMessage) ID() string {
	return "Client.Media.Stop"
}
