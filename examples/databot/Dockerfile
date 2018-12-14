# The build stage
FROM golang:1.11-stretch as builder
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/ann-kilzer/go-keybase-chat-bot/examples/databot
COPY . .
RUN dep ensure -v
RUN go build -o /tmp/databot .

# The final stage
FROM debian:stretch-slim
RUN apt-get update && apt-get install -y curl && \
    curl -o /tmp/keybase.deb https://prerelease.keybase.io/keybase_amd64.deb && \
    dpkg -i /tmp/keybase.deb; apt-get install -f -y && \
    rm /tmp/keybase.deb && apt-get remove -y --purge curl
RUN addgroup --system --gid 200 databot && adduser --system --home /databot --uid 200  --gid 200 databot
USER databot
WORKDIR /databot
RUN mkdir -p downloads
COPY --chown=databot --from=builder /tmp/databot .
VOLUME ["/databot/config", "/databot/downloads"]
CMD ["sh", "-c", "/usr/bin/keybase oneshot && /databot/databot"]
