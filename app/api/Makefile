build:
	[ -d "bin/" ] && rm -r "bin/"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
	mkdir -p bin/
	zip bin/bootstrap.zip bootstrap
	rm bootstrap