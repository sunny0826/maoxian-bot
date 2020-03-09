# build step
FROM golang:1.12 as builder

LABEL maintainer="sunnydog0826@gmail.com"

COPY . /build/

WORKDIR /build

RUN make build

# run step
FROM alpine

ARG BOT_VERSION=0.0.1

RUN apk update \
    && apk add --no-cache ca-certificates \
    && rm -rf /var/cache/apk/*

# copy bin from build step
COPY --from=builder /build/target/mxbot-${BOT_VERSION}/mxbot /bin/

ENTRYPOINT ["/bin/mxbot"]