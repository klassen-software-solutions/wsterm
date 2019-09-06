build:
	go generate
	go build -o wsterm cmd/wsterm/main.go
