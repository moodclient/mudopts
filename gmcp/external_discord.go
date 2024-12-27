package gmcp

import "github.com/moodclient/telnet"

func NewPackageExternalDiscord() Package {
	return Package{
		ID:      "External.Discord",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ExternalDiscordInfoMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[ExternalDiscordStatusMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[ExternalDiscordHelloMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[ExternalDiscordGetMessage],
			},
		},
	}
}

type ExternalDiscordInfoMessage struct {
	BaseMessage

	InviteURL     string `json:"inviteurl,omitempty"`
	ApplicationID string `json:"applicationid,omitempty"`
}

func (m ExternalDiscordInfoMessage) ID() string {
	return "External.Discord.Info"
}

type ExternalDiscordStatusMessage struct {
	BaseMessage

	SmallImage     []string `json:"smallimage"`
	SmallImageText string   `json:"smallimagetext"`
	Details        string   `json:"details"`
	State          string   `json:"state"`
	PartySize      int      `json:"partysize"`
	PartyMax       int      `json:"partymax"`
	Game           string   `json:"game"`
	StartTime      string   `json:"starttime,omitempty"`
	EndTime        string   `json:"endtime,omitempty"`
}

func (m ExternalDiscordStatusMessage) ID() string {
	return "External.Discord.Status"
}

type ExternalDiscordHelloMessage struct {
	BaseMessage

	User    string `json:"user"`
	Private bool   `json:"private"`
}

func (m ExternalDiscordHelloMessage) ID() string {
	return "External.Discord.Hello"
}

type ExternalDiscordGetMessage struct {
	BaseMessage
}

func (m ExternalDiscordGetMessage) ID() string {
	return "External.Discord.Get"
}
