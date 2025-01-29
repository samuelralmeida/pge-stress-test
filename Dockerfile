FROM golang:1.22.10

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o tester

ENTRYPOINT ["./tester"]
