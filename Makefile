build: build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-windows-386 build-windows-amd64 build-darwin-386 build-darwin-amd64

build-linux-386:
	@env GOOS=linux GOARCH=386 go build -o=out/rebyre_linux-386 ./cmd/cli/

build-linux-amd64:
	@env GOOS=linux GOARCH=amd64 go build -o=out/rebyre_linux-amd64 ./cmd/cli/

build-linux-arm:
	@env GOOS=linux GOARCH=arm go build -o=out/rebyre_linux-arm ./cmd/cli/

build-linux-arm64:
	@env GOOS=linux GOARCH=arm64 go build -o=out/rebyre_linux-arm64 ./cmd/cli/
	
build-windows-386:
	@env GOOS=windows GOARCH=386 go build -o=out/rebyre_windows-386.exe ./cmd/cli/
	
build-windows-amd64:
	@env GOOS=linux GOARCH=amd64 go build -o=out/rebyre_windows-amd64.exe ./cmd/cli/

build-darwin-386:
	@env GOOS=darwin GOARCH=386 go build -o=out/rebyre_darwin-386.exe ./cmd/cli/

build-darwin-amd64:
	@env GOOS=darwin GOARCH=amd64 go build -o=out/rebyre_darwin-amd64.exe ./cmd/cli/
