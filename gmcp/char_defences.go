package gmcp

import "github.com/moodclient/telnet"

func NewPackageCharDefences() Package {
	return Package{
		ID:      "Char.Defences",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharDefencesListMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharDefencesAddMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharDefencesRemoveMessage],
			},
		},
	}
}

type Defence struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type CharDefencesListMessage struct {
	BaseMessage

	ValueMessage[[]Defence]
}

func (m CharDefencesListMessage) ID() string {
	return "Char.Defences.List"
}

type CharDefencesAddMessage struct {
	BaseMessage

	Defence
}

func (m CharDefencesAddMessage) ID() string {
	return "Char.Defences.Add"
}

type CharDefencesRemoveMessage struct {
	BaseMessage

	ValueMessage[[]string]
}

func (m CharDefencesRemoveMessage) ID() string {
	return "Char.Defences.Remove"
}
