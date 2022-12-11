FROM alpine

# Using Go 1.10 version

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && apk --no-cache add tzdata
WORKDIR /app
COPY build/p2p-monitor .

# grpc
EXPOSE 8082
# http
EXPOSE 4012


ENTRYPOINT ["./p2p-monitor"]
