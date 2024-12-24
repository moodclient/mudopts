package mudopts

import (
	"fmt"
	"strconv"

	"github.com/moodclient/telnet"
	"github.com/moodclient/telnet/telopts"
)

type ClientCaps uint64

const (
	ANSI ClientCaps = 1 << iota
	VT100
	UTF8
	Colors256
	MouseTracking
	OscColorPalette
	ScreenReader
	Proxy
	TrueColor
	MNES
	MSLP
	SSL
)

type ClientInfo struct {
	Name    string
	Version string

	Charset   string
	IPAddress string
	TermType  string

	Capabilities ClientCaps
}

func (i ClientInfo) RegisterMTTS(usage telnet.TelOptUsage) telnet.TelnetOption {
	return telopts.RegisterTTYPE(usage, []string{
		i.Name, i.TermType, fmt.Sprintf("MTTS %d", int(i.Capabilities)),
	})
}

func (i *ClientInfo) RegisterMNES(usage telnet.TelOptUsage) telnet.TelnetOption {
	vars := make(map[string]string)
	if i.Name != "" {
		vars["CLIENT_NAME"] = i.Name
	}

	if i.Version != "" {
		vars["CLIENT_VERSION"] = i.Version
	}

	if i.IPAddress != "" {
		vars["IPADDRESS"] = i.IPAddress
	}

	vars["MTTS"] = strconv.Itoa(int(i.Capabilities))

	if i.TermType != "" {
		vars["TERMINAL_TYPE"] = i.TermType
	}

	return telopts.RegisterNEWENVIRON(usage, telopts.NEWENVIRONConfig{
		WellKnownVarKeys: []string{"CHARSET", "CLIENT_NAME", "CLIENT_VERSION", "IPADDRESS", "MTTS", "TERMINAL_TYPE"},
		InitialVars:      vars,
	})
}
