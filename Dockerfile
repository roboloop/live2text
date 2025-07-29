FROM golang:1.24.5-alpine

ENV CGO_ENABLED=1

RUN apk add --no-cache build-base pkgconfig portaudio-dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["sh"]
