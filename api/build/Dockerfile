FROM golang:1.24.5

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux && go build -ldflags='-s -w' -buildvcs=false -o ./tmp/api ./api

ENTRYPOINT ["./tmp/api"]