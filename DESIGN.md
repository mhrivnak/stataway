## Main Loop

* Start monitors.
* Listen for triggers on a channel.
* Notify monitors of state changes.

## States

### Home

* If location triggers home, no-op.
* If location triggers away, set state to Away.
* If network triggers home, no-op.

### Away

* If location triggers home, set state to Home.
* If location triggers away, no-op.
* If network triggers home, set state to Home.

## Location

Monitor location with google maps API.

Trigger when threshhold crossed:

* outbound and all devices are outside
* inbound

### TODO

[] Figure out how to get locations.
[] Figure out how to turn locations into a distance from home.

### Device States

* outside
* inside

## Devices on Network

If app state is ``home``, idle.

If app state is ``away``, monitor for presence of devices on the network.

Trigger when any known device is on the network.

### Notification

* Home: stop monitoring
* Away: begin monitoring

## References

### autoaway

https://github.com/MilhouseVH/autoaway.py

A python app that sets thermostat home/away based on device presence on the network.

### Google Maps API

### Venstar API


