package mudopts

import (
	"strings"

	"github.com/moodclient/telnet/telopts"
)

type MCCPCompressionStatusEvent struct {
	telopts.BaseTelOptEvent
	Started bool
	Sending bool
}

func (s MCCPCompressionStatusEvent) String() string {
	var sb strings.Builder
	sb.WriteString(s.Option().String())
	sb.WriteString(": compression status change- compression: ")

	if s.Sending {
		sb.WriteString("SENDING ")
	} else {
		sb.WriteString("RECEIVING ")
	}

	if s.Started {
		sb.WriteString("STARTED")
	} else {
		sb.WriteString("STOPPED")
	}

	return sb.String()
}
