FROM golang:1.13 as builder
WORKDIR /go/src/github.com/thecasualcoder/dobby
COPY ./ ./
RUN make build-deps compile

FROM ubuntu:bionic
COPY --from=builder /go/src/github.com/thecasualcoder/dobby/out/dobby /usr/local/bin/
EXPOSE 4444
RUN useradd -ms /bin/bash dobby
USER dobby
CMD ["dobby", "server", "--bind-address", "0.0.0.0"]
