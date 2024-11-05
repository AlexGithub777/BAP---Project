# Use the official Golang image to build the app
FROM golang:1.22.5 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o edms.exe ./cmd/main.go

# Start a new stage from scratch
FROM gcr.io/distroless/base

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/edms.exe /app/edms.exe

# Copy the .env file into the container
COPY --from=builder /app/.env /app/.env

COPY .env /.env

COPY templates /templates
COPY static /static
COPY internal /internal

# Command to run the executable
CMD ["/app/edms.exe"]
