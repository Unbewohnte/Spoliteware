all: server client
	mkdir -p bin
	mv spolitewareServer* bin
	mv spolitewareClient* bin
	mv thanksdir bin

server:
	cd src/server && CGO_ENABLED=0 go build && mv spolitewareServer ../../ && cp -r thanksdir ../../

client:
	cd src/client && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && mv spolitewareClient ../../spolitewareClient_linux_amd64
	cd src/client && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build && mv spolitewareClient ../../spolitewareClient_darwin_amd64
	cd src/client && CGO_ENABLED=0 GOOS=windows GOARCH=386 go build && mv spolitewareClient.exe ../../spolitewareClient_windows_x32.exe
	cd src/client && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build && mv spolitewareClient.exe ../../spolitewareClient_windows_amd64.exe

release: client server
	mkdir -p release/spoliteware
	cp LICENSE release/spoliteware
	mv spolitewareServer release/spoliteware
	mv thanksdir release/spoliteware
	mv spolitewareClient* release/spoliteware
	cd release && zip -r spoliteware spoliteware