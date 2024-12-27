package gmcp

import (
	"encoding/json"
	"strings"

	"github.com/moodclient/telnet"
)

func NewPackageCharItems() Package {
	return Package{
		ID:      "Char.Items",
		Version: 1,
		Messages: []MessageData{
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CharItemsInvMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CharItemsContentsMessage],
			},
			{
				Sender: telnet.SideClient,
				Create: CreateMessage[CharItemsRoomMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharItemsListMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharItemsAddMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharItemsUpdateMessage],
			},
			{
				Sender: telnet.SideServer,
				Create: CreateMessage[CharItemsRemoveMessage],
			},
		},
	}
}

type CharItemsInvMessage struct {
	BaseMessage
}

func (m CharItemsInvMessage) ID() string {
	return "Char.Items.Inv"
}

type CharItemsContentsMessage struct {
	BaseMessage

	ValueMessage[int]
}

func (m CharItemsContentsMessage) ID() string {
	return "Char.Items.Contents"
}

type CharItemsRoomMessage struct {
	BaseMessage
}

func (m CharItemsRoomMessage) ID() string {
	return "Char.Items.Room"
}

type ItemAttributes int

const (
	ItemAttWorn ItemAttributes = 1 << iota
	ItemAttWearableNotWorn
	ItemAttWieldedLeft
	ItemAttWieldedRight
	ItemAttGroupable
	ItemAttContainer
	ItemAttRiftable
	ItemAttFluid
	ItemAttEditable
	ItemAttMonster
	ItemAttDeadMonster
	ItemAttTakeable
	ItemAttNoTarget
)

func (a ItemAttributes) String() string {
	var sb strings.Builder
	if a&ItemAttWorn != 0 {
		sb.WriteByte('w')
	} else if a&ItemAttWearableNotWorn != 0 {
		sb.WriteByte('W')
	}

	if a&ItemAttWieldedLeft != 0 {
		sb.WriteByte('l')
	}

	if a&ItemAttWieldedRight != 0 {
		sb.WriteByte('L')
	}

	if a&ItemAttGroupable != 0 {
		sb.WriteByte('g')
	}

	if a&ItemAttContainer != 0 {
		sb.WriteByte('c')
	}

	if a&ItemAttRiftable != 0 {
		sb.WriteByte('r')
	}

	if a&ItemAttFluid != 0 {
		sb.WriteByte('f')
	}

	if a&ItemAttEditable != 0 {
		sb.WriteByte('e')
	}

	if a&ItemAttMonster != 0 {
		sb.WriteByte('m')
	}

	if a&ItemAttDeadMonster != 0 {
		sb.WriteByte('d')
	}

	if a&ItemAttTakeable != 0 {
		sb.WriteByte('t')
	}

	if a&ItemAttNoTarget != 0 {
		sb.WriteByte('x')
	}

	return sb.String()
}

func ParseItemAttributes(s string) ItemAttributes {
	var att ItemAttributes

	for i := 0; i < len(s); i++ {
		b := s[i]

		switch b {
		case 'w':
			att |= ItemAttWorn
		case 'W':
			att |= ItemAttWearableNotWorn
		case 'l':
			att |= ItemAttWieldedLeft
		case 'L':
			att |= ItemAttWieldedRight
		case 'g':
			att |= ItemAttGroupable
		case 'c':
			att |= ItemAttContainer
		case 'r':
			att |= ItemAttRiftable
		case 'f':
			att |= ItemAttFluid
		case 'e':
			att |= ItemAttEditable
		case 'm':
			att |= ItemAttMonster
		case 'd':
			att |= ItemAttDeadMonster
		case 't':
			att |= ItemAttTakeable
		case 'x':
			att |= ItemAttNoTarget
		}
	}

	return att
}

func (i *ItemAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *ItemAttributes) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*i = ParseItemAttributes(str)
	return nil
}

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`

	Attrib ItemAttributes `json:"attrib"`
}

type CharItemsListMessage struct {
	BaseMessage

	Location string `json:"location"`
	Items    []Item `json:"items"`
}

func (m CharItemsListMessage) ID() string {
	return "Char.Items.List"
}

type CharItemsAddMessage struct {
	BaseMessage

	Location string `json:"location"`
	Item     Item   `json:"item"`
}

func (m CharItemsAddMessage) ID() string {
	return "Char.Items.Add"
}

type CharItemsUpdateMessage struct {
	BaseMessage

	Location string `json:"location"`
	Item     Item   `json:"item"`
}

func (m CharItemsUpdateMessage) ID() string {
	return "Char.Items.Update"
}

type CharItemsRemoveMessage struct {
	BaseMessage

	Location string `json:"location"`
	Item     Item   `json:"item"`
}

func (m CharItemsRemoveMessage) ID() string {
	return "Char.Items.Remove"
}
