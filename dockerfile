FROM alpine:3.7

EXPOSE 8080

ADD svc-auth /bin/svc-auth

ENTRYPOINT "svc-auth"