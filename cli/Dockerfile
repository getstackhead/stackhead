FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o ./bin/stackhead-cli .

FROM alpine
COPY --from=builder /build/bin/stackhead-cli /bin/
WORKDIR /project
CMD ["/bin/stackhead-cli"]
