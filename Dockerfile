FROM golang:1-alpine3.18 AS builder

RUN go build main.go -o /usr/local/bin/godot-build-tool

FROM alpine:3.18
RUN apk add --no-cache unzip

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.35-r1/glibc-2.35-r1.apk \
    && apk add glibc-2.35-r1.apk

WORKDIR /opt
COPY --from=builder /usr/local/bin/godot-build-tool /opt/godot-build-tool
CMD ["/opt/godot-build-tool"]