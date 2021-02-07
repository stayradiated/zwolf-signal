# zwolf-signal

> a service for interacting with Signal

## Getting Started

You should first use `docker-compose` to build and launch the service.

```
docker-compose up --build --detach
```

Now you will need to verify your signal account. This is currently done
manually.

See the [signal-cli docs](https://github.com/AsamK/signal-cli) for more info.
You will probably need to visit
https://signalcaptchas.org/registration/generate.html to generate a captcha
code.

Your signal session will be stored in the `./signal-cli-config` directory. If
you have issues with messages not being sent or decrypted, try deleting
everything in this directory and reconnecting your account.

```
docker-compose exec main signal-cli -u USERNAME register --captcha CAPTCHA
docker-compose exec main signal-cli -u USERNAME verify CODE --pin PIN
docker-compose exec main signal-cli -u USERNAME send -m "hello world" RECIPIENT
```

You can now restart the app and everything should be good.

```
docker-compose restart main
```

View the logs with

```
docker-compose logs --tail=100 -f main
```

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
