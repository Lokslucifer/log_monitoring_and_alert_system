# Use Go base image
FROM golang:1.23.4-alpine3.19

# Set working directory
WORKDIR /app

# Copy all source code
COPY . .

RUN go mod download

# Build both apps
RUN go build -o log_collector_app ./cmd

# Set back to root app directory

# Run both apps in parallel
CMD ["./log_collector_app"]
