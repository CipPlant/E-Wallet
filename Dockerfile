FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o run -ldflags "-s -w" cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /app/run /app/run
ENTRYPOINT [ "/app/run"]