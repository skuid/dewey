FROM alpine:3.6

LABEL maintainer=ethan@skuid.com

RUN apk add -U ca-certificates
ADD dewey /usr/local/bin/

ENTRYPOINT ["dewey"]