FROM alpine

WORKDIR /app

ADD artifacts/cloud-initer /app/cloud-initer
ADD config.default.json /app/config.json

RUN apk --update upgrade && \
    apk add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

EXPOSE 7002

ENTRYPOINT ["/app/cloud-initer"]