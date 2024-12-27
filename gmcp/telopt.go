package gmcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/moodclient/mudopts"
	"github.com/moodclient/telnet"
	"github.com/moodclient/telnet/telopts"
)

const gmcp telnet.TelOptCode = 201

func RegisterGMCP(usage telnet.TelOptUsage, clientInfo mudopts.ClientInfo, packages ...Package) *GMCP {
	g := &GMCP{
		BaseTelOpt: telopts.NewBaseTelOpt(gmcp, "GMCP", usage),

		clientInfo: clientInfo,
		packages:   make(map[string]Package),

		remoteClientSupported:    make(map[string]int),
		remoteClientIntersection: make(map[string]struct{}),

		clientMessages: make(map[string]MessageFactory),
		serverMessages: make(map[string]MessageFactory),

		messageToPackage: make(map[string]string),
	}

	g.AddPackages(packages...)

	return g
}

type GMCP struct {
	telopts.BaseTelOpt

	parseLock sync.Mutex

	clientInfo mudopts.ClientInfo
	packages   map[string]Package

	remoteClientSupported    map[string]int
	remoteClientIntersection map[string]struct{}

	clientMessages map[string]MessageFactory
	serverMessages map[string]MessageFactory

	messageToPackage map[string]string
}

func (g *GMCP) AddPackages(pkgs ...Package) error {
	g.parseLock.Lock()
	defer g.parseLock.Unlock()

	hadNoPackages := len(g.packages) == 0

	addPackageSet := make(map[string]Package, len(pkgs))
	for _, pkg := range pkgs {
		addPackageSet[pkg.ID] = pkg
	}

	// If we add a different version from what we already support, that is a legitimate
	// change, but otherwise we should ignore the addition
	removePackages := make(map[string]Package)
	for _, pkg := range addPackageSet {
		oldPackage, packageExists := g.packages[pkg.ID]

		if packageExists && oldPackage.Version != pkg.Version {
			removePackages[oldPackage.ID] = oldPackage
		} else if packageExists {
			delete(addPackageSet, pkg.ID)
		}
	}

	if len(removePackages) > 0 {
		g.removePackages(removePackages)
	}

	for _, pkg := range addPackageSet {
		g.packages[pkg.ID] = pkg

		for message := range pkg.AllMessages {
			msg, err := message.Create(g, nil)
			if err != nil {
				return err
			}

			msgID := strings.ToUpper(msg.ID())

			if message.Sender == telnet.SideClient {
				g.clientMessages[msgID] = message.Create
			} else {
				g.serverMessages[msgID] = message.Create
			}

			g.messageToPackage[msgID] = pkg.ID
		}
	}

	if g.Terminal() != nil && g.Terminal().Side() == telnet.SideServer {
		// Update the intersection with client support
		for _, removed := range removePackages {
			clientVersion, clientSupports := g.remoteClientSupported[removed.ID]
			if clientSupports && clientVersion == removed.Version {
				delete(g.remoteClientIntersection, removed.ID)
			}
		}

		for _, added := range addPackageSet {
			clientVersion, clientSupports := g.remoteClientSupported[added.ID]
			if clientSupports && clientVersion == added.Version {
				g.remoteClientIntersection[added.ID] = struct{}{}
			}
		}
	}

	var err error

	if g.Terminal() != nil && g.Terminal().Side() == telnet.SideClient && g.RemoteState() == telnet.TelOptActive {
		// Notify the server of our updated support

		if len(removePackages) > 0 || hadNoPackages {
			// Set support
			msg := CoreSupportsSetMessage{}
			msg.Value = g.packageSupports(g.packages)
			err = g.SendMessage(msg)
		} else {
			// Add support
			msg := CoreSupportsAddMessage{}
			msg.Value = g.packageSupports(addPackageSet)
			err = g.SendMessage(msg)
		}
	}

	return err
}

func (g *GMCP) removePackages(pkgs map[string]Package) error {
	for _, pkg := range pkgs {
		mapPackage, packageExists := g.packages[pkg.ID]
		if !packageExists {
			continue
		}

		delete(g.packages, mapPackage.ID)

		for message := range pkg.AllMessages {
			msg, err := message.Create(g, nil)
			if err != nil {
				return err
			}

			if message.Sender == telnet.SideClient {
				delete(g.clientMessages, msg.ID())
			} else {
				delete(g.serverMessages, msg.ID())
			}
		}
	}

	return nil
}

func (g *GMCP) RemovePackage(pkgs ...Package) error {
	g.parseLock.Lock()
	defer g.parseLock.Unlock()

	pkgRemoveSet := make(map[string]Package)
	for _, pkg := range pkgs {
		pkgRemoveSet[pkg.ID] = pkg
	}

	err := g.removePackages(pkgRemoveSet)
	if err != nil {
		return err
	}

	if g.Terminal() != nil && g.Terminal().Side() == telnet.SideServer {
		// Update the intersection with client support
		for _, removed := range pkgs {
			clientVersion, clientSupports := g.remoteClientSupported[removed.ID]
			if clientSupports && clientVersion == removed.Version {
				delete(g.remoteClientIntersection, removed.ID)
			}
		}
	}

	if g.Terminal() != nil && g.Terminal().Side() == telnet.SideClient && g.RemoteState() == telnet.TelOptActive {
		// Update support
		msg := CoreSupportsRemoveMessage{}
		msg.Value = g.packageSupports(pkgRemoveSet)
		err = g.SendMessage(msg)
	}

	return err
}

func (g *GMCP) packageSupports(pkgs map[string]Package) []string {
	out := make([]string, 0, len(pkgs))

	for _, pkg := range pkgs {
		out = append(out, pkg.Key())
	}

	return out
}

func (g *GMCP) writeMessage(id string, rawJson []byte) error {
	bytes := bytes.NewBuffer(make([]byte, 0, len(id)+len(rawJson)+1))
	bytes.WriteString(id)
	bytes.WriteByte(' ')
	_, err := bytes.Write(rawJson)
	if err != nil {
		return err
	}

	g.Terminal().Keyboard().WriteCommand(telnet.Command{
		OpCode:         telnet.SB,
		Option:         gmcp,
		Subnegotiation: bytes.Bytes(),
	}, nil)

	return nil
}

func (g *GMCP) SendMessage(message Message) error {
	if g.LocalState() != telnet.TelOptActive && g.RemoteState() != telnet.TelOptActive {
		return nil
	}

	// The server shouldn't send messages we know the client doesn't support
	if g.Terminal().Side() == telnet.SideServer {
		pkgID, hasPkg := g.messageToPackage[message.ID()]
		if hasPkg {
			_, clientSupported := g.remoteClientIntersection[pkgID]
			if !clientSupported {
				return nil
			}
		}
	}

	id := message.ID()
	rawJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return g.writeMessage(id, rawJson)
}

func (g *GMCP) TransitionRemoteState(newState telnet.TelOptState) (func() error, error) {
	postFunc, err := g.BaseTelOpt.TransitionRemoteState(newState)
	if err != nil {
		return postFunc, err
	}

	if newState == telnet.TelOptActive && g.Terminal().Side() == telnet.SideClient {
		// Send hello and send support after sending response
		return func() error {
			err = g.SendMessage(CoreHelloMessage{
				Client:  g.clientInfo.Name,
				Version: g.clientInfo.Version,
			})
			if err != nil {
				return err
			}

			return g.SendMessage(CoreSupportsSetMessage{
				ValueMessage: ValueMessage[[]string]{
					Value: g.packageSupports(g.packages),
				},
			})
		}, nil
	}

	return postFunc, err
}

func (g *GMCP) TransitionLocalState(newState telnet.TelOptState) (func() error, error) {
	postFunc, err := g.BaseTelOpt.TransitionLocalState(newState)
	if err != nil {
		return postFunc, err
	}

	if newState == telnet.TelOptInactive && g.Terminal().Side() == telnet.SideServer {
		// Clear client support
		for key := range g.remoteClientSupported {
			delete(g.remoteClientSupported, key)
			delete(g.remoteClientIntersection, key)
		}
	}

	return postFunc, err
}

func (g *GMCP) readMessageName(subnegotiation []byte) (string, int) {
	var sb strings.Builder
	var index int

	for index < len(subnegotiation) {
		r, size := utf8.DecodeRune(subnegotiation[index:])
		index += size

		if r == ' ' {
			break
		}

		sb.WriteRune(r)
	}

	return sb.String(), index
}

func (g *GMCP) createMessage(messageName string, rawJson json.RawMessage) (Message, error) {
	var factory MessageFactory
	var hasFactory bool

	if g.Terminal().Side() == telnet.SideClient {
		// Sender was server
		factory, hasFactory = g.serverMessages[strings.ToUpper(messageName)]
	} else {
		// Sender was client
		factory, hasFactory = g.clientMessages[strings.ToUpper(messageName)]
	}

	var err error
	if !hasFactory {
		msg := UnknownMessage{
			id:         messageName,
			rawMessage: rawJson,
			MapMessage: NewMapMessage(),
		}

		if len(rawJson) > 0 {
			err = json.Unmarshal(rawJson, &msg)
		}

		return msg, err
	}

	return factory(g, rawJson)
}

func (g *GMCP) Subnegotiate(subnegotiation []byte) error {
	if len(subnegotiation) == 0 {
		return errors.New("gmcp: received empty subnegotiation")
	}

	messageName, consumed := g.readMessageName(subnegotiation)
	jsonLen := len(subnegotiation) - consumed

	var rawJson json.RawMessage
	if jsonLen > 0 {
		rawJson = make([]byte, 0, jsonLen)
		rawJson = append(rawJson, subnegotiation[consumed:]...)
	}

	msg, err := g.createMessage(messageName, rawJson)
	if err != nil {
		return err
	}

	g.Terminal().RaiseTelOptEvent(msg)

	return nil
}

func (g *GMCP) SubnegotiationString(subnegotiation []byte) (string, error) {
	if len(subnegotiation) == 0 {
		return "", errors.New("gmcp: received empty subnegotiation")
	}

	messageName, consumed := g.readMessageName(subnegotiation)
	var sb strings.Builder
	sb.WriteString(messageName)

	if consumed < len(subnegotiation) {
		sb.WriteByte(' ')
		_, err := sb.Write(subnegotiation[consumed:])
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}
