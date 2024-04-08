FROM golang:1.18 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o medichat-be .

FROM alpine:latest
COPY --from=builder /app/medichat-be /medichat-be
EXPOSE 8080
CMD ["/medichat-be"]
