FROM alpine:3.14

RUN apk add libc6-compat

COPY trigger /trigger
COPY docker_entrypoint.sh /docker_entrypoint.sh

CMD ["/docker_entrypoint.sh"]
