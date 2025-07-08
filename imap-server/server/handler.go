package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ZiplEix/mail-toolchain/shared/database"
	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

func HandleCommand(session *Session, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid command format")
	}

	tag, cmd := parts[0], strings.ToUpper(parts[1])

	switch cmd {
	case "CAPABILITY":
		session.Send("* CAPABILITY IMAP4rev1 AUTH=PLAIN")
		session.Send(tag + " OK CAPABILITY completed")

	case "NOOP":
		session.Send(tag + " OK NOOP completed")

	case "LOGIN":
		if len(parts) < 4 {
			session.Send(tag + " BAD Usage: LOGIN username password")
			return nil
		}
		session.Username = parts[2]
		session.State = StateAuthenticated
		session.Send(tag + " OK LOGIN completed")

	case "LIST":
		session.Send(`* LIST (\HasNoChildren) "/" "INBOX"`)
		session.Send(tag + " OK LIST completed")

	case "SELECT":
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

	case "FETCH":
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

	case "UID":
		if len(parts) >= 4 && strings.ToUpper(parts[2]) == "FETCH" {
			return handleUIDFetch(session, tag, parts[3], parts[4:])
		}
		session.Send(tag + " BAD Unsuported UID subcommand")
		logger.Info(fmt.Sprintf("Unsuported command: %s %v", tag, parts))

	case "LOGOUT":
		session.Send("* BYE Logging out")
		session.Send(tag + " OK LOGOUT completed")
		session.LoggedOut = true

	default:
		session.Send(tag + " BAD Unknown command")
	}

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
				// Extrait les champs demandés dans les parenthèses
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
				// Corps complet
				bodyLiteral = mail.RawData
				fetchParts = append(fetchParts,
					fmt.Sprintf("BODY[] {%d}", len(mail.RawData)))
				hasBody = true
			}
		}

		if hasBody {
			// 1. FETCH avec {n} (literal)
			session.Send(fmt.Sprintf("* %d FETCH (%s", seq, strings.Join(fetchParts, " ")))
			// 2. Contenu du literal (sur une ligne)
			session.Send(bodyLiteral)
			// 3. Fermeture du FETCH
			session.Send(")")
		} else {
			// Tout sur une ligne car pas de literal
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
