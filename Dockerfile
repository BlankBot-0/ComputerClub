FROM golang:latest as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/task .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/bin/task .
ENTRYPOINT ["./task"]