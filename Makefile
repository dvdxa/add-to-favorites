clean:
	del myapp.exe

build: clean
	go build -tags myapp -o myapp.exe ./cmd/main.go

test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

.PHONY: cover
cover:
	go test -short -count=1  -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
