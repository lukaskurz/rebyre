build: build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-windows-386 build-windows-amd64

build-linux-386:
	@env GOOS=linux GOARCH=386 go build -o=out/rebyre_linux-386 .

build-linux-amd64:
	@env GOOS=linux GOARCH=amd64 go build -o=out/rebyre_linux-amd64 . 

build-linux-arm:
	@env GOOS=linux GOARCH=arm go build -o=out/rebyre_linux-arm .

build-linux-arm64:
	@env GOOS=linux GOARCH=arm64 go build -o=out/rebyre_linux-arm64 .
	
build-windows-386:
	@env GOOS=linux GOARCH=arm64 go build -o=out/rebyre_windows-386.exe .
	
build-windows-amd64:
	@env GOOS=linux GOARCH=arm64 go build -o=out/rebyre_windows-amd64.exe .
	