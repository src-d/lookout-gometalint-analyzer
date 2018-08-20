FROM golang:1.10-alpine

RUN apk add --no-cache git dumb-init
RUN go get -u gopkg.in/alecthomas/gometalinter.v2
ADD ./build/bin/gometalint-analyzer /bin/gometalint-analyzer

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["gometalint-analyzer"]
