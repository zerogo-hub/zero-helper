mod:
	go mod download
	go mod tidy
	go mod verify
	go mod vendor
lint:
	golangci-lint run