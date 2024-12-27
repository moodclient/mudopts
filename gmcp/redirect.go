package gmcp

import "github.com/moodclient/telnet"

func NewPackageRedirect() Package {
	return Package{
		ID:      "Redirect",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[RedirectWindowMessage],
			},
		},
	}
}

type RedirectWindowMessage struct {
	BaseMessage

	ValueMessage[string]
}

func (m RedirectWindowMessage) ID() string {
	return "Redirect.Window"
}
