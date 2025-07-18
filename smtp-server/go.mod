module github.com/ZiplEix/mail-toolchain/smtp-server

go 1.23.1

require (
	github.com/ZiplEix/mail-toolchain/shared v0.0.0-20250703161514-8d3b3016f32e
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/toorop/go-dkim v0.0.0-20250226130143-9025cce95817
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace github.com/ZiplEix/mail-toolchain/shared => ../shared
