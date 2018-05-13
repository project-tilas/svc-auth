FROM alpine:3.7

EXPOSE 8080

RUN apk --no-cache add ca-certificates

ADD svc-auth /bin/svc-auth

ENTRYPOINT "svc-auth"