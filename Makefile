build:
	go build -o bin/odin.exe

test:
	go test -v ./... -count=1
