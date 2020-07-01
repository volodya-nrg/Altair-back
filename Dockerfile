FROM golang:1.13.5 as builder
WORKDIR /go/src/altair
COPY ./app .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -a -v ./cmd/altair/main.go
RUN chmod +x ./main

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY ./app/config-*.json ./
COPY --from=builder /go/src/altair/main ./main

CMD ["/app/main", "-config=config-release.json"]
