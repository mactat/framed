FROM mactat/framed:latest as builder

FROM alpine:latest as tester

RUN apk add --no-cache git bash sudo

# clone bats
RUN git clone https://github.com/bats-core/bats-core.git /test/bats
RUN git clone https://github.com/bats-core/bats-support.git /test/test_helper/bats-support
RUN git clone https://github.com/bats-core/bats-assert.git /test/test_helper/bats-assert
RUN git clone https://github.com/bats-core/bats-file.git /test/test_helper/bats-file

COPY --from=builder /bin/framed /bin/framed
COPY /test/test.bats /test/test.yaml /test/


