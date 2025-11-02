# Dockerfile for Go development environment

# Start from the official Go image.
# Using a specific version ensures consistency.
FROM golang:1.22-alpine

# Set the working directory inside the container.
WORKDIR /app

# Install Air for live-reloading.
# This makes development much faster as the app restarts on code changes.
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum files to download dependencies.
# This step is separated to leverage Docker's layer caching.
# Dependencies are only re-downloaded if these files change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Expose the port the application will run on.
# This is for documentation purposes; the mapping happens in docker-compose.
EXPOSE 8080

# The command to run the application.
# We use 'air' which will watch for file changes and re-compile/re-run the app.
# This is the entry point for the development container.
CMD ["air"]
