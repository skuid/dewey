FROM golang:1.9-alpine as builder

ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/skuid/dewey
ADD . .
RUN go build -o dewey

FROM alpine:3.6
LABEL maintainer=ethan@skuid.com

RUN apk add -U ca-certificates
COPY --from=builder /go/src/github.com/skuid/dewey/dewey /usr/bin/dewey

ENTRYPOINT ["dewey"]