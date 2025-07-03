# SMTP Client (mail-client)

This package provides a reusable SMTP client in Go to send emails via a remote SMTP server.

It supports:
- Custom SMTP configuration (host, port, auth)
- Plain authentication
- Sending plain text emails
- STARTTLS support (on port 587)

---

## 🔧 Installation

Place this package in your project and import it:

```go
import "github.com/ZiplEix/mail-toolchain/mail-client"
```

## 📧 Usage

```go
package main

import (
    "log"

    smtpclient "github.com/yourusername/mail-toolchain/smtp-client"
)

func main() {
    smtpclient.Setup(smtpclient.SMTPConfig{
        Host:     "smtp.gmail.com",
        Port:     587,
        Username: "your@gmail.com",
        Password: "your-app-password", // Use app-specific password
    })

    err := smtpclient.SendMail(
        "your@gmail.com",
        []string{"dest@example.com"},
        "Hello from Go",
        "This is the body of the email.",
    )

    if err != nil {
        log.Fatal("Send failed:", err)
    }

    log.Println("Email sent successfully.")
}
```

## 🛡️ Notes

 - Make sure to enable "Less secure apps" or use an app-specific password if you're using Gmail.
 - Only supports plain text emails for now. HTML and attachments are not supported yet.
 - This client uses Go's standard `net/smtp`.

## 🔮 TODO

 - HTML content support
 - Attachments (MIME)
 - SMTPS support (port 465)
 - DKIM signing
