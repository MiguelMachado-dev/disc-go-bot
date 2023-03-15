.PHONY: build run clean

IMAGE_NAME := golang:1.20
CONTAINER_NAME := disc-go-bot

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME)

clean:
	docker rm -f $(CONTAINER_NAME) || true
	docker rmi -f $(IMAGE_NAME) || true