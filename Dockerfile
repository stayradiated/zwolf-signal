# ==============================================================================
# zwolf-signal
# ==============================================================================

FROM golang:1.15.3-alpine as zwolf-signal

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/stayradiated/zwolf-assistant

COPY go.mod go.sum ./

RUN go mod download

COPY main.go assistant signal ./

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
  && apk add --no-cache dbus \
  # signal-cli
  && ln -s /opt/signal-cli/bin/signal-cli /usr/bin/signal-cli \
  && mkdir -p /home/.local/share/signal-cli

COPY ./dbus/org.asamk.Signal.conf /etc/dbus-1/system.d/
COPY ./dbus/org.asamk.Signal.service /usr/share/dbus-1/system-services/

COPY ./zwolf-signal /usr/bin/zwolf-signal
COPY ./init.sh /root/init.sh

RUN apk add file
RUN file /opt/signal-cli/bin/signal-cli
RUN file /usr/bin/zwolf-signal

COPY --from=zwolf-signal /root/zwolf-signal /usr/bin/zwolf-signal

CMD "/root/init.sh"
# dbus-daemon --system
# signal-cli -u "${USERNAME}" daemon --system
# signal-cli --dbus-system send -m "${MESSAGE}" "${RECIPIENT}"
