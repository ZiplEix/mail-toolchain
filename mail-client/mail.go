package mailclient

import (
	"fmt"
	"strings"
)

func buildMessage(from string, to []string, subject, body string) string {
	header := map[string]string{
		"From":    from,
		"To":      strings.Join(to, ", "),
		"Subject": subject,
	}

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body + "\r\n")

	return msg.String()
}
