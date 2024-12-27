package gmcp

import "github.com/moodclient/telnet"

func NewPackageCore() Package {
	return Package{
		ID:      "Core",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CoreHelloMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CoreSupportsSetMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CoreSupportsAddMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CoreSupportsRemoveMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CoreKeepAliveMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CorePingClientMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CorePingServerMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CoreGoodbyeMessage],
			},
		},
	}
}

type CoreHelloMessage struct {
	BaseMessage

	Client  string `json:"client"`
	Version string `json:"version"`
}

func (m CoreHelloMessage) ID() string {
	return "Core.Hello"
}

type CoreSupportsSetMessage struct {
	BaseMessage

	ValueMessage[[]string]
}

func (m CoreSupportsSetMessage) ID() string {
	return "Core.Supports.Set"
}

type CoreSupportsAddMessage struct {
	BaseMessage

	ValueMessage[[]string]
}

func (m CoreSupportsAddMessage) ID() string {
	return "Core.Supports.Add"
}

type CoreSupportsRemoveMessage struct {
	BaseMessage

	ValueMessage[[]string]
}

func (m CoreSupportsRemoveMessage) ID() string {
	return "Core.Supports.Remove"
}

type CoreKeepAliveMessage struct {
	BaseMessage
}

func (m CoreKeepAliveMessage) ID() string {
	return "Core.KeepAlive"
}

type CorePingClientMessage struct {
	BaseMessage

	ValueMessage[int]
}

func (m CorePingClientMessage) ID() string {
	return "Core.Ping"
}

type CorePingServerMessage struct {
	BaseMessage
}

func (m CorePingServerMessage) ID() string {
	return "Core.Ping"
}

type CoreGoodbyeMessage struct {
	BaseMessage

	ValueMessage[string]
}

func (m CoreGoodbyeMessage) ID() string {
	return "Core.Goodbye"
}
