build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap endpoints/hello/main.go
	mkdir -p bin/
	zip bin/bootstrap.zip bootstrap
	rm bootstrap