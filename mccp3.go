package mudopts

import (
	"compress/zlib"
	"io"

	"github.com/moodclient/telnet"
	"github.com/moodclient/telnet/telopts"
)

const (
	mccp3 telnet.TelOptCode = 87
)

func RegisterMCCP3(usage telnet.TelOptUsage) *MCCP3 {
	return &MCCP3{
		BaseTelOpt: telopts.NewBaseTelOpt(mccp3, "MCCP3", usage),
	}
}

type MCCP3 struct {
	telopts.BaseTelOpt
	compressed bool
}

func (m *MCCP3) TransitionRemoteState(newState telnet.TelOptState) (func() error, error) {
	afterFunc, err := m.BaseTelOpt.TransitionRemoteState(newState)
	if err != nil {
		return afterFunc, err
	}

	if newState == telnet.TelOptActive {
		return func() error {
			m.Terminal().Keyboard().WriteCommand(telnet.Command{
				OpCode: telnet.SB,
				Option: mccp3,
			}, func() error {
				// Start sending compression
				err := m.Terminal().Keyboard().WrapWriter(func(writer io.Writer) (io.Writer, error) {
					return zlib.NewWriter(writer), nil
				})
				if err != nil {
					return err
				}

				// Send event
				m.compressed = true
				m.Terminal().RaiseTelOptEvent(MCCPCompressionStatusEvent{
					BaseTelOptEvent: telopts.BaseTelOptEvent{m},
					Started:         true,
					Sending:         true,
				})
				return nil
			})
			return nil
		}, nil
	}

	if newState == telnet.TelOptInactive {
		return func() error {
			// Stop sending compression
			err := m.Terminal().Keyboard().WrapWriter(func(writer io.Writer) (io.Writer, error) {
				return writer, nil
			})
			if err != nil {
				return err
			}

			// Send event
			m.compressed = false
			m.Terminal().RaiseTelOptEvent(MCCPCompressionStatusEvent{
				BaseTelOptEvent: telopts.BaseTelOptEvent{m},
				Started:         false,
				Sending:         true,
			})

			return nil
		}, nil
	}

	return afterFunc, err
}

func (m *MCCP3) TransitionLocalState(newState telnet.TelOptState) (func() error, error) {
	afterFunc, err := m.BaseTelOpt.TransitionLocalState(newState)
	if err != nil {
		return afterFunc, err
	}

	if newState == telnet.TelOptInactive && m.compressed {
		// Stop receiving compression
		err := m.Terminal().Printer().WrapReader(func(reader io.Reader) (io.Reader, error) {
			return reader, nil
		})
		if err != nil {
			return afterFunc, err
		}

		// Send event
		m.compressed = false
		m.Terminal().RaiseTelOptEvent(MCCPCompressionStatusEvent{
			BaseTelOptEvent: telopts.BaseTelOptEvent{m},
			Started:         false,
			Sending:         false,
		})
	}

	return afterFunc, err
}

func (m *MCCP3) Subnegotiate(subnegotiation []byte) error {
	if m.LocalState() == telnet.TelOptActive {
		// Start receiving compression
		err := m.Terminal().Printer().WrapReader(func(reader io.Reader) (io.Reader, error) {
			return zlib.NewReader(reader)
		})
		if err != nil {
			return err
		}

		// Send event
		m.compressed = true
		m.Terminal().RaiseTelOptEvent(MCCPCompressionStatusEvent{
			BaseTelOptEvent: telopts.BaseTelOptEvent{m},
			Started:         true,
			Sending:         false,
		})
	}

	return m.BaseTelOpt.Subnegotiate(subnegotiation)
}

func (m *MCCP3) SubnegotiationString(subnegotiation []byte) (string, error) {
	return "BEGIN COMPRESSION", nil
}
