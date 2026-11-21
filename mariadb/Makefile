include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=raspberrypi5 port=$(MARIA_DB_PORT) dbname=$(MARIA_DB_NAME) user=$(MARIA_DB_USER) password=$(MARIA_DB_PPASSWRD)"

install-deps:
	set GOBIN=$(LOCAL_BIN) && go install github.com/pressly/goose/v3/cmd/goose@latest
