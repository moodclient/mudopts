package gmcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/moodclient/telnet"
	"github.com/moodclient/telnet/telopts"
)

type Message interface {
	ID() string
	RawMessage() json.RawMessage

	telnet.TelOptEvent
}

type MessageInitialize interface {
	InitializeAsEvent(g *GMCP, message Message, raw json.RawMessage)
}

type BaseMessage struct {
	telopts.BaseTelOptEvent `json:"-"`

	idCache    string          `json:"-"`
	rawMessage json.RawMessage `json:"-"`
}

func (m BaseMessage) RawMessage() json.RawMessage {
	return m.rawMessage
}

func (m *BaseMessage) InitializeAsEvent(g *GMCP, message Message, raw json.RawMessage) {
	m.rawMessage = raw
	m.idCache = message.ID()
	m.BaseTelOptEvent = telopts.BaseTelOptEvent{TelnetOption: g}
}

func (m BaseMessage) String() string {
	var sb strings.Builder
	sb.WriteString(m.Option().String())
	sb.WriteByte(':')
	sb.WriteByte(' ')

	if m.idCache == "" {
		sb.WriteString("Outbound Event")
	} else {
		sb.WriteString(m.idCache)
		sb.WriteByte(' ')
		sb.WriteByte('-')
		sb.WriteByte(' ')
		sb.WriteString(string(m.rawMessage))
	}

	return sb.String()
}

type MessageFactory func(g *GMCP, raw json.RawMessage) (Message, error)

type MessageData struct {
	Sender telnet.TerminalSide
	Create MessageFactory
}

type Package struct {
	ID       string
	Version  int
	Messages []MessageData
}

func (p Package) AllMessages(yield func(MessageData) bool) {
	for _, message := range p.Messages {
		if !yield(message) {
			return
		}
	}
}

func (p Package) Key() string {
	return fmt.Sprintf("%s %d", p.ID, p.Version)
}

type CreateMessageConstraint[T Message] interface {
	*T
	MessageInitialize
}

func CreateMessage[T Message, U CreateMessageConstraint[T]](g *GMCP, raw json.RawMessage) (Message, error) {
	var zero T

	var ptr = U(&zero)

	err := InitializeMessage[T, U](g, raw, ptr)
	return zero, err
}

func InitializeMessage[T Message, U CreateMessageConstraint[T]](g *GMCP, raw json.RawMessage, message U) error {
	if message == nil {
		return fmt.Errorf("initializemessage: cannot initialize nil message")
	}

	if len(raw) > 0 {
		err := json.Unmarshal([]byte(raw), message)
		if err != nil {
			return err
		}
	}

	message.InitializeAsEvent(g, *message, raw)

	return nil
}

type ValueMessage[T any] struct {
	Value T
}

func (m ValueMessage[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func (m *ValueMessage[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Value)
}

type MapMessage struct {
	ignoreKeys map[string]struct{}
	values     map[string]string
}

func NewMapMessage(ignoreKeys ...string) MapMessage {
	msg := MapMessage{
		ignoreKeys: make(map[string]struct{}),
		values:     make(map[string]string),
	}
	for _, key := range ignoreKeys {
		msg.ignoreKeys[key] = struct{}{}
	}

	return msg
}

func (m MapMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.values)
}

func (m *MapMessage) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, m.values)
}

func (m *MapMessage) Value(key string) (string, bool) {
	value, exists := m.values[key]
	return value, exists
}

func (m *MapMessage) SetValue(key, value string) {
	m.values[key] = value
}

func (m *MapMessage) Keys(yield func(string) bool) {
	for key := range m.values {
		_, ignored := m.ignoreKeys[key]

		if !ignored && !yield(key) {
			return
		}
	}
}

type UnknownMessage struct {
	telopts.BaseTelOptEvent

	id         string
	rawMessage json.RawMessage
	MapMessage
}

func (m UnknownMessage) String() string {
	var sb strings.Builder
	sb.WriteString("Unknown GMCP Message: ")
	sb.WriteString(m.id)
	sb.WriteByte(' ')
	sb.WriteByte('-')
	sb.WriteByte(' ')
	sb.WriteString(string(m.rawMessage))

	return sb.String()
}

func (m UnknownMessage) ID() string {
	return m.id
}

func (m UnknownMessage) RawMessage() json.RawMessage {
	return m.rawMessage
}
