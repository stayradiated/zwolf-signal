# ==============================================================================
# zwolf-signal
# ==============================================================================

FROM golang:1.15.3-alpine as zwolf-signal

WORKDIR $GOPATH/src/github.com/stayradiated/zwolf-assistant

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./
COPY assistant ./assistant
COPY signal ./signal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/zwolf-signal

# ==============================================================================
# SIGNAL-CLI
# ==============================================================================

FROM alpine:3.12.1 as signal-cli

ARG SIGNAL_CLI_VERSION="0.6.11"

RUN \
  wget \
    -O /tmp/signal-cli.tgz \
    "https://github.com/AsamK/signal-cli/releases/download/v${SIGNAL_CLI_VERSION}/signal-cli-${SIGNAL_CLI_VERSION}.tar.gz"
RUN tar xzvf /tmp/signal-cli.tgz -C /tmp
RUN mv "/tmp/signal-cli-${SIGNAL_CLI_VERSION}" /opt/signal-cli

# ==============================================================================
# RELEASE
# ==============================================================================

FROM openjdk:16-alpine

COPY --from=signal-cli /opt/signal-cli /opt/signal-cli

RUN true \
  apk update \
  && apk add --no-cache dbus dbus-x11 \
  && ln -s /opt/signal-cli/bin/signal-cli /usr/bin/signal-cli \
  && mkdir -p /home/.local/share/signal-cli

COPY --from=zwolf-signal /go/bin/zwolf-signal /usr/bin/zwolf-signal

COPY ./init.sh /init.sh
CMD "/init.sh"
