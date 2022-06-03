mod:
	go mod download
	go mod tidy
	go mod verify
	go mod vendor
lint:
	golangci-lint run
test:
	go test ./... -coverprofile cover.out && go tool cover -html=cover.out -o cover.html