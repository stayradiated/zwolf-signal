#!/bin/sh

export DBUS_SESSION_BUS_ADDRESS=$(dbus-daemon --session --fork --print-address)
/usr/bin/zwolf-signal
