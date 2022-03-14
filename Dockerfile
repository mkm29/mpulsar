# builder image
FROM golang:1.17.7-alpine3.15 as builder
ENV LOGLEVEL=DEBUG
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o gopulsar ./cmd/mpulsar/main.go

# generate clean, final image for end users
FROM alpine:3.15.0
COPY --from=builder /build/gopulsar .

EXPOSE 4000

CMD ["./gopulsar"]