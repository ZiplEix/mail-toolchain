package server

import (
	"fmt"
	"strings"
)

type imapHandler func(session *Session, tag string, parts []string) error

type commandEntry struct {
	Handler  imapHandler
	NeedAuth bool
}

var imapCommands = map[string]commandEntry{
	"CAPABILITY":   {capability, false},
	"AUTHENTICATE": {authenticate, false},
	"LOGIN":        {login, false},
	"NOOP":         {noop, false},
	"LIST":         {list, true},
	"SELECT":       {selectMailbox, true},
	"FETCH":        {fetch, true},
	"UID":          {uid, true},
	"LOGOUT":       {logout, false},
}

func HandleCommand(session *Session, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid command format")
	}

	tag, cmd := parts[0], strings.ToUpper(parts[1])

	entry, ok := imapCommands[cmd]
	if !ok {
		session.Send(tag + " BAD Unknown command: " + cmd)
		return nil
	}

	if entry.NeedAuth && session.State == StateNotAuthenticated {
		session.Send(tag + " NO Authentication required for command: " + cmd)
		return nil
	}

	return entry.Handler(session, tag, parts)
}
