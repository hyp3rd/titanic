# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13.0-buster as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
# Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix 'static' -o titanic cmd/titanic/main.go

# Start a new stage from scratch
FROM alpine:3.10.2

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/titanic .

EXPOSE 3000 8443

# Command to run the executable
CMD ["./titanic"]

# Metadata
LABEL org.opencontainers.image.vendor="Hyperd" \
    org.opencontainers.image.url="https://hyperd.sh" \
    org.opencontainers.image.title="Titanic Development Image" \
    org.opencontainers.image.description="Container Solution API-exercise" \
    org.opencontainers.image.version="v0.5" \
    org.opencontainers.image.documentation="https://gitlab.com/hyperd/titanic/blob/master/README.md"
