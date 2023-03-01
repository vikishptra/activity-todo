# Use a golang alpine image as the base image
FROM golang:1.18.2-alpine3.15 AS builder

# Install Git
RUN apk update && apk add --no-cache git

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Clean up any existing go.mod and go.sum files, and create a new one
RUN go mod tidy

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=v1.0.0'" .

# Start from a new scratch image
FROM scratch

# Set environment variables for MySQL
ENV MYSQL_HOST=172.17.0.1
ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=root
ENV MYSQL_DBNAME=todo4
ENV MYSQL_PORT=3306

ENV REDIS_HOST=172.17.0.1
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=

# Expose port 3030
EXPOSE 3030

# Copy the binary from the builder stage
COPY --from=builder /app/vikishptra /usr/bin/

# Set the working directory to /app
WORKDIR /app

# Copy the config file
COPY --from=builder /app/config.json /app

# Set the entrypoint and default command
ENTRYPOINT ["vikishptra", "todos"]
CMD ["--v"]