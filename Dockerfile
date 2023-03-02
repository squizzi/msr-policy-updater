FROM golang:1.19.5-alpine3.17 AS builder

RUN apk add make bash git openssh build-base curl
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
RUN mkdir -p /go/src
ADD . /go/src/
WORKDIR /go/src

RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
RUN --mount=type=ssh go build -o ./bin/linux/msr-policy-updater main.go

FROM alpine:latest AS msr-policy-updater
COPY --from=builder /go/src/bin/linux/msr-policy-updater /

ENTRYPOINT [ "/msr-policy-updater" ]
