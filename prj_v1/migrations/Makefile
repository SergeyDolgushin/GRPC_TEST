include .env

LOCAL_BIN:=F:/GO/GRPC_TEST/prj_v1/bin

LOCAL_MIGRATION_DIR=$(CURDIR)
# LOCAL_MIGRATION_DSN="host=raspberrypi5 port=$(MARIA_DB_PORT) dbname=$(MARIA_DB_NAME) user=$(MARIA_DB_USER) password=$(MARIA_DB_PPASSWRD)"
LOCAL_MIGRATION_DSN="$(MARIA_DB_USER):$(MARIA_DB_PASSWRD)@tcp(raspberrypi5:$(MARIA_DB_PORT))/$(MARIA_DB_NAME)"

# dsn := "testuser:testpass@tcp(raspberrypi5:3306)/testdb"
install-deps:
	set GOBIN=$(LOCAL_BIN) && go install github.com/pressly/goose/v3/cmd/goose@latest


local-migration-status:
	$(LOCAL_BIN)/goose.exe -dir $(LOCAL_MIGRATION_DIR) mysql $(LOCAL_MIGRATION_DSN) status -v

local-migration-up:
	$(LOCAL_BIN)/goose.exe -dir $(LOCAL_MIGRATION_DIR) mysql $(LOCAL_MIGRATION_DSN) up -v

local-migration-down:
	$(LOCAL_BIN)/goose.exe -dir $(LOCAL_MIGRATION_DIR) mysql $(LOCAL_MIGRATION_DSN) down -v