package gmcp

import "github.com/moodclient/telnet"

func NewPackageCharSkills() Package {
	return Package{
		ID:      "Char.Skills",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CharSkillsGetMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharSkillsGroupsMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharSkillsListMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharSkillsInfoMessage],
			},
		},
	}
}

type CharSkillsGetMessage struct {
	BaseMessage

	Group string `json:"group"`
	Name  string `json:"name"`
}

func (m CharSkillsGetMessage) ID() string {
	return "Char.Skills.Get"
}

type CharSkillsGroupsMessage struct {
	BaseMessage

	Name string `json:"name"`
	Rank string `json:"rank"`
}

func (m CharSkillsGroupsMessage) ID() string {
	return "Char.Skills.Groups"
}

type CharSkillsListMessage struct {
	BaseMessage

	Group       string   `json:"group"`
	Description []string `json:"desc"`
	List        []string `json:"list"`
}

func (m CharSkillsListMessage) ID() string {
	return "Char.Skills.List"
}

type CharSkillsInfoMessage struct {
	BaseMessage

	Group string `json:"group"`
	Skill string `json:"skill"`
	Info  string `json:"info"`
}

func (m CharSkillsInfoMessage) ID() string {
	return "Char.Skills.Info"
}
