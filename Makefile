build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o m3u8-downloader.exe ./main.go
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o m3u8-downloader ./main.go
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o m3u8-downloader main.go