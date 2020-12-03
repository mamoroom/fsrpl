BINARY_NAME=fsrpl

.PHONY: build lint test test-local emulator realease

clean: 
	rm -fr ./bin/*
build: clean
	go build -o ./bin/$(BINARY_NAME) ./...

lint:
	golangci-lint run ./...

test-local:
	go test ./...

test: lint test-local
	FIRESTORE_EMULATOR_HOST=0.0.0.0:8080 go test ./... -tags=integration

restore2emulator:
	FIRESTORE_EMULATOR_HOST=localhost:8080 go run ./... restore "importData/*" --path "./" --emulators-project-id ${PROJECT}

emulator-down:
	cd emulator && docker-compose down && cd ..
	unset FIRESTORE_EMULATOR_HOST


release:
	goreleaser --rm-dist

