package gmcp

import (
	"github.com/moodclient/telnet"
)

func NewPackageCharAfflictions() Package {
	return Package{
		ID:      "Char.Afflictions",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharAfflictionsListMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharAfflictionsAddMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharAfflictionsRemoveMessage],
			},
		},
	}
}

type Affliction struct {
	Name        string `json:"name"`
	Cure        string `json:"cure"`
	Description string `json:"desc"`
}

type CharAfflictionsListMessage struct {
	BaseMessage

	ValueMessage[[]Affliction]
}

func (m CharAfflictionsListMessage) ID() string {
	return "Char.Afflictions.List"
}

type CharAfflictionsAddMessage struct {
	BaseMessage

	ValueMessage[[]Affliction]
}

func (m CharAfflictionsAddMessage) ID() string {
	return "Char.Afflictions.Add"
}

type CharAfflictionsRemoveMessage struct {
	BaseMessage

	ValueMessage[[]string]
}

func (m CharAfflictionsRemoveMessage) ID() string {
	return "Char.Afflictions.Remove"
}
