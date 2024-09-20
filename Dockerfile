FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install github.com/cespare/reflex@latest

RUN apk add --no-cache git
RUN apk add --no-cache postgresql-dev
RUN apk add --no-cache postgresql-dev gcc musl-dev
RUN go install github.com/cosmtrek/air@v1.27.3

COPY . .
COPY .env .env

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
