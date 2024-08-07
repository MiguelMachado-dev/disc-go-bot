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

# Build the Docker image
docker-build:
	docker build -t $(APP_NAME):$(IMAGE_TAG) -f $(DOCKERFILE_PATH) .

# Run the application in a Docker container
docker-run: docker-build
	docker run -d -p $(PORT):$(PORT) --name $(APP_NAME) --restart always $(APP_NAME):$(IMAGE_TAG)

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

.PHONY: all build docker-build docker-run docker-stop docker-restart clean test run
