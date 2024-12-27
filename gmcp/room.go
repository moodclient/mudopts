package gmcp

import "github.com/moodclient/telnet"

func NewPackageRoom() Package {
	return Package{
		ID:      "Room",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[RoomInfoMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[RoomWrongDirMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[RoomPlayersMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[RoomAddPlayerMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[RoomRemovePlayerMessage],
			},
		},
	}
}

type RoomInfoMessage struct {
	BaseMessage

	Number      int            `json:"num"`
	Name        string         `json:"name"`
	Area        string         `json:"area"`
	Environment string         `json:"environment"`
	Coords      string         `json:"coords"`
	Map         string         `json:"map"`
	Details     []string       `json:"details"`
	Exits       map[string]int `json:"exits"`
}

func (m RoomInfoMessage) ID() string {
	return "Room.Info"
}

type RoomWrongDirMessage struct {
	BaseMessage

	ValueMessage[string]
}

func (m RoomWrongDirMessage) ID() string {
	return "Room.WrongDir"
}

type RoomPlayersMessage struct {
	BaseMessage

	ValueMessage[[]CharName]
}

func (m RoomPlayersMessage) ID() string {
	return "Room.Players"
}

type RoomAddPlayerMessage struct {
	BaseMessage

	CharName
}

func (m RoomAddPlayerMessage) ID() string {
	return "Room.AddPlayer"
}

type RoomRemovePlayerMessage struct {
	BaseMessage

	ValueMessage[string]
}

func (m RoomRemovePlayerMessage) ID() string {
	return "Room.RemovePlayer"
}
