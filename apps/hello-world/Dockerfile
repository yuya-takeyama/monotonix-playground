FROM golang:1.23.5 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM gcr.io/distroless/static-debian12:530158861eebdbbf149f7e7e67bfe45eb433a35c

COPY --from=builder /app/main /
EXPOSE 8080

CMD ["/main"]
