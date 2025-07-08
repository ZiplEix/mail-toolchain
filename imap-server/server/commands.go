package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

func noop(session *Session, tag string, parts []string) error {
	session.Send(tag + " OK NOOP completed")
	return nil
}

func capability(session *Session, tag string, parts []string) error {
	session.Send("* CAPABILITY IMAP4rev1 AUTH=PLAIN")
	session.Send(tag + " OK CAPABILITY completed")

	return nil
}

func authenticate(session *Session, tag string, parts []string) error {
	if len(parts) < 3 {
		session.Send(tag + " BAD Usage: AUTHENTICATE <mechanism>")
		return nil
	}

	mech := strings.ToUpper(parts[2])
	if mech != "PLAIN" {
		session.Send(tag + " BAD Unsupported authentication mechanism: " + mech)
		return nil
	}

	session.Send("+")

	line, err := session.ReadLine()
	if err != nil {
		session.Send(tag + " NO Failed to read authentication string")
		return err
	}

	data, err := base64.StdEncoding.DecodeString(line)
	if err != nil {
		session.Send(tag + " NO Invalid base64 encoding")
	}

	dataParts := bytes.Split(data, []byte{0})
	if len(dataParts) != 3 {
		session.Send(tag + " NO Invalid PLAIN authentication format")
		return nil
	}

	// authzID := string(dataParts[0])
	username := string(dataParts[1])
	password := string(dataParts[2])

	if valid, err := database.CheckUserPassword(username, password); err != nil || !valid {
		session.Send(tag + " NO Authentication failed")
		return nil
	}

	session.Username = username
	session.State = StateAuthenticated
	session.Send(tag + " OK AUTHENTICATE completed")

	return nil
}

func login(session *Session, tag string, parts []string) error {
	if len(parts) < 4 {
		session.Send(tag + " BAD Usage: LOGIN username password")
		return nil
	}
	username := parts[2]
	password := parts[3]

	valid, err := database.CheckUserPassword(username, password)
	if err != nil || !valid {
		session.Send(tag + " NO LOGIN failed")
		return nil
	}

	session.Username = username
	session.State = StateAuthenticated
	session.Send(tag + " OK LOGIN completed")

	return nil
}

func list(session *Session, tag string, parts []string) error {
	session.Send(`* LIST (\HasNoChildren) "/" "INBOX"`)
	session.Send(tag + " OK LIST completed")

	return nil
}

func selectMailbox(session *Session, tag string, parts []string) error {
	session.SelectedMB = "INBOX"
	session.State = StateSelected

	mails, err := database.GetAllMails()
	if err != nil {
		session.Send(tag + " NO Unable to access mailbox")
		return fmt.Errorf("failed to select mailbox: %w", err)
	}

	count := len(mails)

	session.Send(fmt.Sprintf("* %d EXISTS", count))
	session.Send(tag + " OK [READ-WRITE] SELECT completed")

	return nil
}

func fetch(session *Session, tag string, parts []string) error {
	if session.State != StateSelected {
		session.Send(tag + " BAD SELECT a mailbox first")
		return nil
	}

	mails, err := database.GetAllMails()
	if err != nil {
		session.Send(tag + " NO Could not retrieve messages")
		return fmt.Errorf("failed to fetch mails: %w", err)
	}

	fmt.Println("Fetched", len(mails), "mails from database")
	for _, mail := range mails {
		fmt.Println("Mail ID:", mail.ID, "Sender:", mail.Sender, "Recipients:", mail.Recipients, "Received At:", mail.ReceivedAt)
	}

	if len(mails) == 0 {
		session.Send("* 0 EXISTS")
		session.Send(tag + " OK FETCH completed (no messages)")
		return nil
	}

	for i, mail := range mails {
		seq := i + 1 // Why the fuck does IMAP sequence start at 1 ??????
		body := mail.RawData
		size := len(body)

		session.Send(fmt.Sprintf(`* %d FETCH (FLAGS (\Seen) INTERNALDATE "%s" RFC822.SIZE %d BODY[] {%d}`,
			seq, mail.ReceivedAt.Format("02-Jan-2006 15:04:05 -0700"), size, size))

		session.Send(body)
		session.Send(")")
	}

	session.Send(tag + " OK FETCH completed")

	return nil
}

func uid(session *Session, tag string, parts []string) error {
	if len(parts) >= 4 && strings.ToUpper(parts[2]) == "FETCH" {
		return handleUIDFetch(session, tag, parts[3], parts[4:])
	}
	session.Send(tag + " BAD Unsuported UID subcommand")
	logger.Info(fmt.Sprintf("Unsuported command: %s %v", tag, parts))

	return nil
}

func logout(session *Session, tag string, parts []string) error {
	session.Send("* BYE Logging out")
	session.Send(tag + " OK LOGOUT completed")
	session.LoggedOut = true

	return nil
}

type UIDFetchRequest struct {
	RangeStart int
	RangeEnd   int
	Fields     []string
}

func handleUIDFetch(session *Session, tag string, rangeStr string, fields []string) error {
	req, err := parseUIDFetch(rangeStr, fields)
	if err != nil {
		session.Send(tag + " BAD Invalid UID FETCH syntax: " + err.Error())
		return err
	}

	mails, err := database.GetMailsInUIDRange(req.RangeStart, req.RangeEnd)
	if err != nil {
		session.Send(tag + " NO Could not retrieve messages in UID range")
		return fmt.Errorf("db ftech failed: %w", err)
	}

	for i, mail := range mails {
		seq := i + 1
		var fetchParts []string
		var bodyLiteral string
		var hasBody bool

		fetchParts = append(fetchParts, fmt.Sprintf("UID %d", mail.ID))

		for _, field := range req.Fields {
			switch {
			case strings.EqualFold(field, "UID"):
				fetchParts = append(fetchParts, fmt.Sprintf("UID %d", mail.ID))

			case strings.EqualFold(field, "FLAGS"):
				fetchParts = append(fetchParts, "FLAGS (\\Seen)")

			case strings.EqualFold(field, "RFC822.SIZE"):
				fetchParts = append(fetchParts, fmt.Sprintf("RFC822.SIZE %d", len(mail.RawData)))

			case strings.HasPrefix(strings.ToUpper(field), "BODY.PEEK[HEADER.FIELDS"):
				start := strings.Index(field, "(")
				end := strings.LastIndex(field, ")")
				if start == -1 || end == -1 || end <= start {
					session.Send(tag + " BAD malformed HEADER.FIELDS")
					return nil
				}

				fieldsRaw := field[start+1 : end]
				fieldList := strings.Fields(fieldsRaw)
				headerData := extractHeaders(mail.RawData, fieldList)

				bodyLiteral = headerData
				fetchParts = append(fetchParts,
					fmt.Sprintf("BODY[HEADER.FIELDS (%s)] {%d}", strings.Join(fieldList, " "), len(headerData)))
				hasBody = true

			case strings.EqualFold(field, "BODY[]"):
				bodyLiteral = mail.RawData
				fetchParts = append(fetchParts,
					fmt.Sprintf("BODY[] {%d}", len(mail.RawData)))
				hasBody = true
			}
		}

		if hasBody {
			session.Send(fmt.Sprintf("* %d FETCH (%s", seq, strings.Join(fetchParts, " ")))
			session.Send(bodyLiteral)
			session.Send(")")
		} else {
			session.Send(fmt.Sprintf("* %d FETCH (%s)", seq, strings.Join(fetchParts, " ")))
		}
	}

	session.Send(tag + " OK UID FETCH completed")

	return nil
}

func parseUIDFetch(rangeStr string, fields []string) (*UIDFetchRequest, error) {
	var start, end int
	var err error

	if rangeStr == "*" {
		start = -1
		end = -1
	} else if strings.Contains(rangeStr, ":") {
		parts := strings.Split(rangeStr, ":")
		start, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid range start: %s", parts[0])
		}
		if parts[1] == "*" {
			end = -1
		} else {
			end, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end: %s", parts[1])
			}
		}
	} else {
		start, err = strconv.Atoi(rangeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid range: %s", rangeStr)
		}
		end = start
	}

	joined := strings.Join(fields, " ")
	joined = strings.TrimPrefix(joined, "(")
	joined = strings.TrimSuffix(joined, ")")

	fieldList := strings.Fields(joined)

	return &UIDFetchRequest{
		RangeStart: start,
		RangeEnd:   end,
		Fields:     fieldList,
	}, nil
}

func extractHeaders(raw string, wanted []string) string {
	lines := strings.Split(raw, "\r\n")
	var builder strings.Builder
	headerMap := make(map[string]string)

	for _, line := range lines {
		if line == "" {
			break
		}

		colon := strings.Index(line, ":")
		if colon > 0 {
			key := strings.TrimSpace(line[:colon])
			value := strings.TrimSpace(line[colon+1:])
			headerMap[strings.ToLower(key)] = value
		}
	}

	for _, field := range wanted {
		v, ok := headerMap[strings.ToLower(field)]
		if ok {
			builder.WriteString(fmt.Sprintf("%s: %s\r\n", field, v))
		}
	}

	builder.WriteString("\r\n")
	return builder.String()
}
