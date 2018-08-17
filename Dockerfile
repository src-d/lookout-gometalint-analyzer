FROM debian:stretch-slim

RUN apt-get update && apt-get install -y dumb-init

ADD ./build/bin/gometalint-analyzer /bin/gometalint-analyzer

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["gometalint-analyzer"]