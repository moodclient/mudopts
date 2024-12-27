package gmcp

import "github.com/moodclient/telnet"

func NewPackageCharLogin() Package {
	return Package{
		ID:      "Char.Login",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharLoginDefaultMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharLoginResultMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CharLoginCredentialsMessage],
			},
		},
	}
}

type CharLoginDefaultMessage struct {
	BaseMessage

	Type     []string `json:"type"`
	Location string   `json:"location"`
}

func (m CharLoginDefaultMessage) ID() string {
	return "Char.Login.Default"
}

type CharLoginResultMessage struct {
	BaseMessage

	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (m CharLoginResultMessage) ID() string {
	return "Char.Login.Result"
}

type CharLoginCredentialsMessage struct {
	BaseMessage

	Account  string `json:"account"`
	Password string `json:"password"`
}

func (m CharLoginCredentialsMessage) ID() string {
	return "Char.Login.Credentials"
}
