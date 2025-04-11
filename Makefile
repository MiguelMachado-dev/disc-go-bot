# Variables
APP_NAME = disc-go-bot
IMAGE_TAG = latest
DOCKERFILE_PATH = Dockerfile
PORT = 8080

# Default target
all: build

# Build the application
build:
	go build -o $(APP_NAME)

# Build as Windows GUI application (no console window)
build-gui:
	go build -ldflags="-H=windowsgui" -o $(APP_NAME)

# Cross compile for Windows GUI application
windows-gui:
	GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o $(APP_NAME).exe

# Build the Docker image
docker-build:
	docker build \
		--build-arg DISCORD_BOT_TOKEN=$$(grep DISCORD_BOT_TOKEN .env | cut -d '=' -f2) \
		--build-arg COMMANDS_CHANNEL_ID=$$(grep COMMANDS_CHANNEL_ID .env | cut -d '=' -f2) \
		--build-arg TWITCH_CLIENT_ID=$$(grep TWITCH_CLIENT_ID .env | cut -d '=' -f2) \
		--build-arg TWITCH_CLIENT_SECRET=$$(grep TWITCH_CLIENT_SECRET .env | cut -d '=' -f2) \
		--build-arg ENCRYPTION_KEY=$$(grep ENCRYPTION_KEY .env | cut -d '=' -f2) \
		-t migtito/$(APP_NAME):$(IMAGE_TAG) -f $(DOCKERFILE_PATH) .

# Run the application in a Docker container
docker-run: docker-build
	docker run -d -p $(PORT):$(PORT) --name $(APP_NAME) --restart always migtito/$(APP_NAME):$(IMAGE_TAG)

# Stop and remove the running Docker container
docker-stop:
	docker stop $(APP_NAME)

# Restart running Docker container
docker-restart:
	docker-stop
	docker-run

# Clean up build artifacts
clean:
	rm -f $(APP_NAME)

# Execute tests
test:
	go test -v ./...

# Build and run the application locally
run: build
	./$(APP_NAME)

# Build as GUI and run the application locally
run-gui: build-gui
	./$(APP_NAME)

.PHONY: all build build-gui windows-gui docker-build docker-run docker-stop docker-restart clean test run run-gui
