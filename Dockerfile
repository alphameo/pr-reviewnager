FROM golang:1.25.4-alpine AS builder

RUN apk --no-cache add git make curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN set -ex && \
    wget -O /tmp/migrate.tgz https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz && \
    tar -xzf /tmp/migrate.tgz -C /tmp && \
    mv /tmp/migrate /usr/local/bin/migrate && \
    rm /tmp/migrate.tgz # <--- Исправлено: migrate, а не migrate.linux-amd64

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pr-reviewnager ./cmd/pr-reviewnager

FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    wget -O /tmp/migrate.tgz https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz && \
    tar -xzf /tmp/migrate.tgz -C /tmp && \
    mv /tmp/migrate /usr/local/bin/migrate && \
    rm /tmp/migrate.tgz

COPY --from=builder /app/pr-reviewnager /pr-reviewnager

EXPOSE 8080

CMD ["/pr-reviewnager"]
