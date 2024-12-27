package gmcp

import "github.com/moodclient/telnet"

func NewPackageCommChannel() Package {
	return Package{
		ID:      "Comm.Channel",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CommChannelPlayersClientMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CommChannelEnableMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CommChannelPlayersServerMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CommChannelListMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CommChannelTextMessage],
			},
		},
	}
}

type CommChannelPlayersClientMessage struct {
	BaseMessage
}

func (m CommChannelPlayersClientMessage) ID() string {
	return "Comm.Channel.Players"
}

type CommChannelEnableMessage struct {
	BaseMessage

	ValueMessage[string]
}

func (m CommChannelEnableMessage) ID() string {
	return "Comm.Channel.Enable"
}

type CommPlayer struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels,omitempty"`
}

type CommChannelPlayersServerMessage struct {
	BaseMessage

	ValueMessage[[]CommPlayer]
}

func (m CommChannelPlayersServerMessage) ID() string {
	return "Comm.Channel.Players"
}

type CommChannel struct {
	Name    string `json:"name"`
	Caption string `json:"caption"`
	Command string `json:"command"`
}

type CommChannelListMessage struct {
	BaseMessage

	ValueMessage[[]CommChannel]
}

func (m CommChannelListMessage) ID() string {
	return "Comm.Channel.List"
}

type CommChannelTextMessage struct {
	BaseMessage

	Channel string `json:"channel"`
	Talker  string `json:"talker"`
	Text    string `json:"text"`
}

func (m CommChannelTextMessage) ID() string {
	return "Comm.Channel.Text"
}
