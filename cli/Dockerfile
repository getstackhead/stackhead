FROM golang:alpine as builder
WORKDIR /build
ADD . /build
RUN go build -o ./bin/stackhead-cli .

FROM alpine
COPY --from=builder /build/bin/stackhead-cli /bin/
WORKDIR /project
CMD ["/bin/stackhead-cli"]
