# Step 1 ✅

Basic SMTP server for reception of emails
  - Receive emails via SMTP
  - Store emails in a mails/ folder

# Step 2 ✅

credible SMTP server:
  1. ✅ multiple recipients
  2. ✅ address validation
  3. ✅ logs
  4. ✅ error gestion
  5. ✅ RSET, NOOP and VRFY support
  6. ✅ STARTTLS support
  7. ✅ store mails in postgresql

goal: functional and robust SMTP receptor server usable by a true client (Thunderbird, Outlook, etc.).

# 🔁 Step 3 — Sending Emails (SMTP Client)

Now, the sending side.
 - Create an SMTP client in Go (net/smtp) to:
 - Connect to a remote server (Gmail, Outlook, etc.)
 - Send an email with Subject, Body, To, From fields

📦 Goal: Be able to send an email using your own toolchain.

# 📬 Step 4 — Minimal IMAP Server

Future mail client to be able to read incoming messages.
 - TCP server on port 143 (or a custom one)
 - Basic authentication (user/pass or none at first)
 - Core commands: LOGIN, LIST, SELECT, FETCH, LOGOUT
 - Read maisl from db

📂 Goal: Create my own webmail or use an IMAP client to read your emails.

# 🔒 Step 5 — Authentication & Security

 - 🔑 SMTP Auth (LOGIN / PLAIN)
 - 🛡️ STARTTLS (TLS encryption)
 - 🧾 SPF / DKIM / DMARC (required so your emails aren’t marked as spam)
 - 🧅 DNS configuration: MX record, DKIM record, SPF TXT, etc.

📦 Goal: Ensure your emails land in the inbox, not the spam folder.

# 🌐 Step 6 — Webmail / Mobile App / REST API

Create a small frontend or REST API to:

 - View the inbox
 - Send an email
 - Delete / sort messages

Do this with SvelteKit.

# 🧪 Step 7 — Real Interoperability

Connect a real email client to the server (Thunderbird, Outlook, etc.)
 - Send myself an email from Gmail and see it arrive
 - Reply to an email via my SMTP and check it gets delivered properly