# zwolf-signal

> a service for interacting with Signal

## Under the hood

- Uses [`AsamK/signal-cli`](https://github.com/AsamK/signal-cli) to send and
  receive messages

## Todo

- Remove assistant code
- Connect to amqp watermill
- Remove os.Getenv call from `./signal/signal.go`
- Add hooks for calling register/verify/trust

## Thanks

Started as a fork of
[vberezny/signal-assistant](https://github.com/vberezny/signal-assistant) which
made setting up the dbus connection possible.
