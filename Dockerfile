FROM golang:1.25.4-alpine AS builder

RUN apk --no-cache add git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pr-reviewnager ./cmd/pr-reviewnager

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/pr-reviewnager /pr-reviewnager

EXPOSE 8080

CMD ["/pr-reviewnager"]
