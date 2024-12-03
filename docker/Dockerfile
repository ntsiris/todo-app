# Start with the official Golang image as the build stage
FROM golang:1.22 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o todo-app cmd/main.go

# Use a minimal base image for the runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/todo-app .

# Copy any required assets like .env file (if applicable)
# Uncomment the line below if an .env file is required at runtime
# COPY .env .env

# Expose the application's default port
EXPOSE 8080

# Command to run the application
CMD ["./todo-app"]
