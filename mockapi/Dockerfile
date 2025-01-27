# backend/Dockerfile (with air for live reload)
FROM golang:1.23.5-alpine

WORKDIR /app

# Copy go.mod and go.sum files
#COPY go.mod go.sum ./
COPY go.mod ./

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Download dependencies
RUN go mod download

# Copy the rest of the application code
#COPY . .

RUN export PATH="$PATH:$GOBIN"

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application with air
CMD air --build.cmd "go build -o tmp/main ." --build.bin "./tmp/main" --build.exclude_dir "templates,build" --build.poll "true"
