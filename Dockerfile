FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apt-get -y install postgresql-client
RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go
FROM alpine:latest
COPY --from=builder main /app/main
EXPOSE 8080
CMD ["/app/main"]
