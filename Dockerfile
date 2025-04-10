# Use the official Golang image as the base image
FROM golang:1.20

# Build args for configuration
ARG DISCORD_BOT_TOKEN
ARG COMMANDS_CHANNEL_ID
ARG TWITCH_CLIENT_ID
ARG TWITCH_CLIENT_SECRET

# Set as environment variables
ENV DISCORD_BOT_TOKEN=$DISCORD_BOT_TOKEN
ENV COMMANDS_CHANNEL_ID=$COMMANDS_CHANNEL_ID
ENV TWITCH_CLIENT_ID=$TWITCH_CLIENT_ID
ENV TWITCH_CLIENT_SECRET=$TWITCH_CLIENT_SECRET

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files into the container
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Expose the application's port
EXPOSE 8080

# Start the application
CMD ["./main"]
