.PHONY: build run build-server build-client run-server run-client

build: build-server build-client

build-server:
	go build -o server cmd/server/main.go

build-client:
	go build -o client cmd/client/main.go

run-server: build-server
	./server

run-client: build-client
	./client

run: 
	@./server & \
	SERVER_PID=$$!; \
	echo "Server PID: $$SERVER_PID"; \
	sleep 1; \
	./client; \
	STATUS=$$?; \
	echo "Killing server..."; \
	kill $$SERVER_PID; \
	wait $$SERVER_PID 2>/dev/null || true; \
	exit $$STATUS
