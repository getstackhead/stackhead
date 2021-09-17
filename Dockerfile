FROM golang:alpine as builder
WORKDIR /build
COPY . /build
RUN sh /build/.build/build.sh

FROM pad92/ansible-alpine:2.10.3
COPY --from=builder /build/bin/stackhead-cli /bin/
WORKDIR /project
CMD ["/bin/stackhead-cli"]
