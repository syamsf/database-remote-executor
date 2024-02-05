build-web:
	cd cmd/web && GOOS=linux GOARCH=amd64 go build -o ../../bin/web

build-web-mac:
	cd cmd/web && go build -o ../../bin/web-mac

run-web:
	@go run cmd/web/main.go
