# Mail Toolchain

This project is a full self-hosted mail toolchain written in Go.

It includes:

- ✅ A working SMTP server for receiving emails
- 🛠️ A planned SMTP client for sending emails
- 📨 An IMAP server for reading emails
- 🧰 Optional webmail and REST API
- 🔒 Auth, TLS, SPF/DKIM/DMARC for secure delivery

All components are built from scratch using the standard Go libraries, without relying on third-party full mail servers.

### Structure

- `smtp-server/` → Custom SMTP server
- `imap-server/` → (Planned) IMAP server
- `smtp-client/` → (Planned) Go SMTP client
- `webmail/` → (Planned) SvelteKit frontend
- `plan.md` → Detailed roadmap

### Requirements

- Go ≥ 1.20
- Docker (for PostgreSQL)
- Thunderbird / Outlook (for testing)
- OpenSSL (for TLS testing)

### Project Plan

The full project roadmap is available in [./plan.md](./plan.md)
