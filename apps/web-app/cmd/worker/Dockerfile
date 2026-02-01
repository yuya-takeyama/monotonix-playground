FROM golang:1.25.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/worker

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/main /

CMD ["/main"]
