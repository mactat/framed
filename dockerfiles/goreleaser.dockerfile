ARG ALPINE_VERSION=3.18

FROM alpine:${ALPINE_VERSION} as release
COPY framed /bin/framed
CMD [ "/bin/sh" ]