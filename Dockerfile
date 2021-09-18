FROM golang:alpine as builder
WORKDIR /build
COPY . /build
RUN sh /build/.build/build.sh

FROM alpine:3.13.6
COPY --from=builder /build/bin/stackhead-cli /bin/
WORKDIR /project
CMD ["/bin/stackhead-cli"]
