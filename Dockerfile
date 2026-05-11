FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

COPY . .

RUN go build -o main .

CMD ["tail", "-f", "/dev/null"]