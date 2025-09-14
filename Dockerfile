FROM golang:1.24.5

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go clean -cache

ENTRYPOINT ["air"]