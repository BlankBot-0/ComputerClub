FROM golang:latest as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go test -c ./... -o bin/

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/bin/* .
ADD ./test/test_files/* /root/test/test_files/
ADD ./test/test_files/outputs/* /root/test/test_files/outputs/
ENTRYPOINT ["./src.test"]