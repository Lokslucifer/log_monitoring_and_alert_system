FROM golang:1.23.4-alpine3.19

WORKDIR /app

# Install GCC and build tools for CGO
RUN apk add --no-cache build-base sqlite-dev

COPY . .

RUN go mod download

# Enable CGO explicitly
ENV CGO_ENABLED=1

# Build your app
RUN go build -o ./log_processor_app ./cmd

CMD ["./log_processor_app"]
