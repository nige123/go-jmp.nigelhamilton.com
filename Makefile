APP_NAME := jmp
CMD_PATH := ./cmd/jmp
DIST_DIR := dist

.PHONY: build test tidy lint clean cross

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(APP_NAME) $(CMD_PATH)

test:
	go test ./...

tidy:
	go mod tidy

lint:
	go test ./...

clean:
	rm -rf $(APP_NAME) $(DIST_DIR)

cross:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(DIST_DIR)/$(APP_NAME)-linux-386 $(CMD_PATH)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe $(CMD_PATH)
