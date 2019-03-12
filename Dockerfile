FROM golang:1.11.4 as builder
WORKDIR /go/src/github.com/thecasualcoder/dobby
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY ./ ./
RUN make build-deps compile

FROM ubuntu:bionic
COPY --from=builder /go/src/github.com/thecasualcoder/dobby/out/dobby /usr/local/bin/
CMD ["dobby", "server"]
