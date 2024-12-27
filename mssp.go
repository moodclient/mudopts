package mudopts

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/moodclient/telnet"
	"github.com/moodclient/telnet/telopts"
)

const mssp telnet.TelOptCode = 70

const (
	msspVAR byte = 1
	msspVAL byte = 2
)

type MSSPData struct {
	Name    string
	Players int
	Uptime  time.Time

	Charset    []string
	Codebase   []string
	Contact    string
	CrawlDelay int
	Created    int
	DiscordURL *url.URL
	Hostname   string
	Icon       *url.URL
	IP         string
	IPV6       string
	Language   string
	Location   string
	MinimumAge int
	Port       []int
	Referral   []string
	SSLPort    int
	Website    *url.URL

	Family     []string
	Genre      string
	Gameplay   string
	Status     string
	GameSystem string
	InterMUD   []string
	Subgenre   string

	Areas     int
	HelpFiles int
	Mobiles   int
	Objects   int
	Rooms     int
	Classes   int
	Levels    int
	Races     int
	Skills    int

	ANSI           bool
	UTF8           bool
	VT100          bool
	XTerm256Color  bool
	XTermTrueColor bool

	PayToPlay   bool
	PayForPerks bool

	HiringBuilders bool
	HiringCoders   bool
}

type MSSPUpdatedEvent struct {
	telopts.BaseTelOptEvent

	Data MSSPData
}

func (e MSSPUpdatedEvent) String() string {
	return "MSSP Data Updated"
}

type MSSP struct {
	telopts.BaseTelOpt

	data atomic.Pointer[MSSPData]
}

func RegisterMSSP(usage telnet.TelOptUsage, data MSSPData) *MSSP {
	opt := &MSSP{
		BaseTelOpt: telopts.NewBaseTelOpt(mssp, "MSSP", usage),
	}
	opt.data.Store(&data)

	return opt
}

func (m *MSSP) TransitionLocalState(newState telnet.TelOptState) (func() error, error) {
	postFunc, err := m.BaseTelOpt.TransitionLocalState(newState)
	if err != nil {
		return postFunc, err
	}

	if newState == telnet.TelOptActive {
		return func() error {
			buffer := bytes.NewBuffer(nil)
			data := *m.data.Load()
			m.writeToBuffer(buffer, data)

			m.Terminal().Keyboard().WriteCommand(telnet.Command{
				OpCode:         telnet.SB,
				Option:         mssp,
				Subnegotiation: buffer.Bytes(),
			}, nil)

			return nil
		}, nil
	}

	return postFunc, nil
}

func (m *MSSP) Subnegotiate(subnegotiation []byte) error {
	var data MSSPData

	m.readFromBuffer(subnegotiation, &data)

	m.data.Store(&data)
	m.Terminal().RaiseTelOptEvent(MSSPUpdatedEvent{
		BaseTelOptEvent: telopts.BaseTelOptEvent{m},
		Data:            data,
	})

	return nil
}

func (m *MSSP) SubnegotiationString(subnegotiation []byte) (string, error) {
	var index int
	var sb strings.Builder

	for index < len(subnegotiation) && subnegotiation[index] == msspVAR {
		sb.WriteString("VAR ")
		index++

		token, consumed := m.readToken(subnegotiation[index:])
		sb.WriteByte('\'')
		sb.WriteString(token)
		sb.WriteByte('\'')
		sb.WriteByte(' ')
		index += consumed

		for index < len(subnegotiation) && subnegotiation[index] == msspVAL {
			sb.WriteString("VAL ")
			index++

			token, consumed := m.readToken(subnegotiation[index:])
			sb.WriteByte('\'')
			sb.WriteString(token)
			sb.WriteByte('\'')
			sb.WriteByte(' ')
			index += consumed
		}
	}

	return sb.String(), nil
}

func (m *MSSP) Data() MSSPData {
	ptr := m.data.Load()
	if ptr == nil {
		return MSSPData{}
	}

	return *ptr
}

func (m *MSSP) readToken(b []byte) (string, int) {
	var sb strings.Builder

	for i := 0; i < len(b); i++ {
		if b[i] == msspVAL || b[i] == msspVAR {
			return sb.String(), i
		}

		sb.WriteByte(b[i])
	}

	return sb.String(), len(b)
}

func (m *MSSP) readInt(b []byte) (val int, consumed int) {
	token, consumed := m.readToken(b)

	val, err := strconv.Atoi(token)
	if err != nil {
		return 0, consumed
	}

	return val, consumed
}

func (m *MSSP) readBool(b []byte) (val bool, consumed int) {
	num, consumed := m.readInt(b)
	return num > 0, consumed
}

func (m *MSSP) readTime(b []byte) (val time.Time, consumed int) {
	num, consumed := m.readInt(b)
	return time.Unix(int64(num), 0), consumed
}

func (m *MSSP) readURL(b []byte) (val *url.URL, consumed int) {
	token, consumed := m.readToken(b)
	url, err := url.Parse(token)
	if err != nil {
		return nil, consumed
	}

	return url, consumed
}

func msspReadSlice[T comparable](b []byte, readMethod func([]byte) (T, int)) ([]T, int) {
	var slice []T
	var zero T

	var index int

	for index < len(b) && b[index] == msspVAL {
		index++

		val, consumed := readMethod(b[index:])
		if val != zero {
			slice = append(slice, val)
		}
		index += consumed
	}

	return slice, index
}

func msspLastValue[T any](slice []T) T {
	var zero T
	if len(slice) > 0 {
		return slice[len(slice)-1]
	}

	return zero
}

func (m *MSSP) readFromBuffer(b []byte, data *MSSPData) {
	var index int

	for index < len(b) && b[index] == msspVAR {
		index++

		token, consumed := m.readToken(b[index:])
		index += consumed

		var strSlice []string
		var intSlice []int
		var timeSlice []time.Time
		var urlSlice []*url.URL
		var boolSlice []bool

		switch strings.ToUpper(token) {
		case "NAME":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Name = msspLastValue(strSlice)

		case "PLAYERS":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Players = msspLastValue(intSlice)

		case "UPTIME":
			timeSlice, consumed = msspReadSlice(b[index:], m.readTime)
			data.Uptime = msspLastValue(timeSlice)

		case "CHARSET":
			data.Charset, consumed = msspReadSlice(b[index:], m.readToken)

		case "CODEBASE":
			data.Codebase, consumed = msspReadSlice(b[index:], m.readToken)

		case "CONTACT":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Contact = msspLastValue(strSlice)

		case "CRAWL DELAY":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.CrawlDelay = msspLastValue(intSlice)

		case "CREATED":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Created = msspLastValue(intSlice)

		case "DISCORD":
			urlSlice, consumed = msspReadSlice(b[index:], m.readURL)
			data.DiscordURL = msspLastValue(urlSlice)

		case "HOSTNAME":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Hostname = msspLastValue(strSlice)

		case "ICON":
			urlSlice, consumed = msspReadSlice(b[index:], m.readURL)
			data.Icon = msspLastValue(urlSlice)

		case "IP":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.IP = msspLastValue(strSlice)

		case "IPV6":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.IPV6 = msspLastValue(strSlice)

		case "LANGUAGE":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Language = msspLastValue(strSlice)

		case "LOCATION":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Location = msspLastValue(strSlice)

		case "MINIMUM AGE":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.MinimumAge = msspLastValue(intSlice)

		case "PORT":
			data.Port, consumed = msspReadSlice(b[index:], m.readInt)

		case "REFERRAL":
			data.Referral, consumed = msspReadSlice(b[index:], m.readToken)

		case "SSL":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.SSLPort = msspLastValue(intSlice)

		case "WEBSITE":
			urlSlice, consumed = msspReadSlice(b[index:], m.readURL)
			data.Website = msspLastValue(urlSlice)

		case "FAMILY":
			data.Family, consumed = msspReadSlice(b[index:], m.readToken)

		case "GENRE":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Genre = msspLastValue(strSlice)

		case "GAMEPLAY":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Gameplay = msspLastValue(strSlice)

		case "STATUS":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Status = msspLastValue(strSlice)

		case "GAMESYSTEM":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.GameSystem = msspLastValue(strSlice)

		case "INTERMUD":
			data.InterMUD, consumed = msspReadSlice(b[index:], m.readToken)

		case "SUBGENRE":
			strSlice, consumed = msspReadSlice(b[index:], m.readToken)
			data.Subgenre = msspLastValue(strSlice)

		case "AREAS":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Areas = msspLastValue(intSlice)

		case "HELPFILES":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.HelpFiles = msspLastValue(intSlice)

		case "MOBILES":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Mobiles = msspLastValue(intSlice)

		case "OBJECTS":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Objects = msspLastValue(intSlice)

		case "ROOMS":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Rooms = msspLastValue(intSlice)

		case "CLASSES":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Classes = msspLastValue(intSlice)

		case "LEVELS":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Levels = msspLastValue(intSlice)

		case "RACES":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Races = msspLastValue(intSlice)

		case "SKILLS":
			intSlice, consumed = msspReadSlice(b[index:], m.readInt)
			data.Skills = msspLastValue(intSlice)

		case "ANSI":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.ANSI = msspLastValue(boolSlice)

		case "UTF-8":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.UTF8 = msspLastValue(boolSlice)

		case "VT100":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.VT100 = msspLastValue(boolSlice)

		case "XTERM 256 COLORS":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.XTerm256Color = msspLastValue(boolSlice)

		case "XTERM TRUE COLORS":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.XTermTrueColor = msspLastValue(boolSlice)

		case "PAY TO PLAY":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.PayToPlay = msspLastValue(boolSlice)

		case "PAY FOR PERKS":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.PayForPerks = msspLastValue(boolSlice)

		case "HIRING BUILDERS":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.HiringBuilders = msspLastValue(boolSlice)

		case "HIRING CODERS":
			boolSlice, consumed = msspReadSlice(b[index:], m.readBool)
			data.HiringCoders = msspLastValue(boolSlice)
		}

		index += consumed
	}
}

func (m *MSSP) writeString(buffer *bytes.Buffer, name string, value string) {
	buffer.WriteByte(msspVAR)
	buffer.WriteString(name)
	buffer.WriteByte(msspVAL)
	buffer.WriteString(value)
}

func (m *MSSP) writeStringSlice(buffer *bytes.Buffer, name string, value []string) {
	buffer.WriteByte(msspVAR)
	buffer.WriteString(name)

	for _, str := range value {
		buffer.WriteByte(msspVAL)
		buffer.WriteString(str)
	}
}

func (m *MSSP) writeInt(buffer *bytes.Buffer, name string, value int) {
	m.writeString(buffer, name, strconv.Itoa(value))
}

func (m *MSSP) writeIntSlice(buffer *bytes.Buffer, name string, value []int) {
	buffer.WriteByte(msspVAR)
	buffer.WriteString(name)

	for _, val := range value {
		buffer.WriteByte(msspVAL)
		buffer.WriteString(strconv.Itoa(val))
	}
}

func (m *MSSP) writeBool(buffer *bytes.Buffer, name string, value bool) {
	str := "0"
	if value {
		str = "1"
	}

	m.writeString(buffer, name, str)
}

func (m *MSSP) writeTime(buffer *bytes.Buffer, name string, value time.Time) {
	m.writeString(buffer, name, strconv.Itoa(int(value.Unix())))
}

func (m *MSSP) writeURL(buffer *bytes.Buffer, name string, url *url.URL) {
	if url == nil {
		m.writeString(buffer, name, "")
		return
	}

	m.writeString(buffer, name, url.String())
}

func (m *MSSP) writeToBuffer(buffer *bytes.Buffer, data MSSPData) {
	m.writeString(buffer, "NAME", data.Name)
	m.writeInt(buffer, "PLAYERS", data.Players)
	m.writeTime(buffer, "UPTIME", data.Uptime)

	if len(data.Charset) > 0 {
		m.writeStringSlice(buffer, "CHARSET", data.Charset)
	}

	if len(data.Codebase) > 0 {
		m.writeStringSlice(buffer, "CODEBASE", data.Codebase)
	}

	if data.Contact != "" {
		m.writeString(buffer, "CONTACT", data.Contact)
	}

	if data.CrawlDelay != 0 {
		m.writeInt(buffer, "CRAWL DELAY", data.CrawlDelay)
	}

	if data.Created != 0 {
		m.writeInt(buffer, "CREATED", data.Created)
	}

	if data.DiscordURL != nil {
		m.writeURL(buffer, "DISCORD", data.DiscordURL)
	}

	if data.Hostname != "" {
		m.writeString(buffer, "HOSTNAME", data.Hostname)
	}

	if data.Icon != nil {
		m.writeURL(buffer, "ICON", data.Icon)
	}

	if data.IP != "" {
		m.writeString(buffer, "IP", data.IP)
	}

	if data.IPV6 != "" {
		m.writeString(buffer, "IPV6", data.IPV6)
	}

	if data.Language != "" {
		m.writeString(buffer, "LANGUAGE", data.Language)
	}

	if data.Location != "" {
		m.writeString(buffer, "LOCATION", data.Location)
	}

	if data.MinimumAge > 0 {
		m.writeInt(buffer, "MINIMUM AGE", data.MinimumAge)
	}

	if len(data.Port) > 0 {
		m.writeIntSlice(buffer, "PORT", data.Port)
	}

	if len(data.Referral) > 0 {
		m.writeStringSlice(buffer, "REFERRAL", data.Referral)
	}

	if data.SSLPort > 0 {
		m.writeInt(buffer, "SSL", data.SSLPort)
	}

	if data.Website != nil {
		m.writeURL(buffer, "WEBSITE", data.Website)
	}

	if len(data.Family) > 0 {
		m.writeStringSlice(buffer, "FAMILY", data.Family)
	}

	if data.Genre != "" {
		m.writeString(buffer, "GENRE", data.Genre)
	}

	if data.Gameplay != "" {
		m.writeString(buffer, "GAMEPLAY", data.Gameplay)
	}

	if data.Status != "" {
		m.writeString(buffer, "STATUS", data.Status)
	}

	if data.GameSystem != "" {
		m.writeString(buffer, "GAMESYSTEM", data.GameSystem)
	}

	if len(data.InterMUD) > 0 {
		m.writeStringSlice(buffer, "INTERMUD", data.InterMUD)
	}

	if data.Subgenre != "" {
		m.writeString(buffer, "SUBGENRE", data.Subgenre)
	}

	if data.Areas > 0 {
		m.writeInt(buffer, "AREAS", data.Areas)
	}

	if data.HelpFiles > 0 {
		m.writeInt(buffer, "HELPFILES", data.HelpFiles)
	}

	if data.Mobiles > 0 {
		m.writeInt(buffer, "MOBILES", data.Mobiles)
	}

	if data.Objects > 0 {
		m.writeInt(buffer, "OBJECTS", data.Objects)
	}

	if data.Rooms > 0 {
		m.writeInt(buffer, "ROOMS", data.Rooms)
	}

	if data.Classes > 0 {
		m.writeInt(buffer, "CLASSES", data.Classes)
	}

	if data.Levels > 0 {
		m.writeInt(buffer, "LEVELS", data.Levels)
	}

	if data.Races > 0 {
		m.writeInt(buffer, "RACES", data.Races)
	}

	if data.Skills > 0 {
		m.writeInt(buffer, "SKILLS", data.Skills)
	}

	m.writeBool(buffer, "ANSI", data.ANSI)
	m.writeBool(buffer, "UTF-8", data.UTF8)
	m.writeBool(buffer, "VT100", data.VT100)
	m.writeBool(buffer, "XTERM 256 COLORS", data.XTerm256Color)
	m.writeBool(buffer, "XTERM TRUE COLORS", data.XTermTrueColor)

	m.writeBool(buffer, "PAY TO PLAY", data.PayToPlay)
	m.writeBool(buffer, "PAY FOR PERKS", data.PayForPerks)

	m.writeBool(buffer, "HIRING BUILDERS", data.HiringBuilders)
	m.writeBool(buffer, "HIRING CODERS", data.HiringCoders)
}
