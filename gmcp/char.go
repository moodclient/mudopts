package gmcp

import (
	"encoding/json"

	"github.com/moodclient/telnet"
)

func NewPackageChar() Package {
	return Package{
		ID:      "Char",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CharLoginMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharNameMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: func(g *GMCP, raw json.RawMessage) (Message, error) {
					msg := CharVitalsMessage{
						MapMessage: NewMapMessage("string"),
					}

					err := InitializeMessage(g, raw, &msg)
					return msg, err
				},
			},
			{
				Sender: telnet.SideServer,
				Create: func(g *GMCP, raw json.RawMessage) (Message, error) {
					m := CharStatusVarsMessage{
						MapMessage: NewMapMessage(),
					}

					err := InitializeMessage(g, raw, &m)
					return m, err
				},
			},
			{
				Sender: telnet.SideServer,
				Create: func(g *GMCP, raw json.RawMessage) (Message, error) {
					m := CharStatusMessage{
						MapMessage: NewMapMessage(),
					}

					err := InitializeMessage(g, raw, &m)
					return m, err
				},
			},
		},
	}
}

type CharLoginMessage struct {
	BaseMessage

	Name     string `json:"name"`
	Password string `json:"password"`
}

func (m CharLoginMessage) ID() string {
	return "Char.Login"
}

type CharName struct {
	Name     string `json:"name"`
	FullName string `json:"fullname"`
}

type CharNameMessage struct {
	BaseMessage

	CharName
}

func (m CharNameMessage) ID() string {
	return "Char.Name"
}

type CharVitalsMessage struct {
	BaseMessage

	MapMessage
}

func (m CharVitalsMessage) ID() string {
	return "Char.Vitals"
}

func (m *CharVitalsMessage) StringValue() string {
	val, _ := m.MapMessage.Value("string")
	return val
}

type CharStatusVarsMessage struct {
	BaseMessage

	MapMessage
}

func (m CharStatusVarsMessage) ID() string {
	return "Char.StatusVars"
}

type CharStatusMessage struct {
	BaseMessage

	MapMessage
}

func (m CharStatusMessage) ID() string {
	return "Char.Status"
}
