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

### locationsharinglib

https://github.com/costastf/locationsharinglib

A python library that logs into Google and gets location sharing data by
scraping endpoints that appear intended for browser-based use.

### Google Maps API

There doesn't appear to be a public API for accessing location sharing info.
Please update this or file an issue if you find one!

### Venstar API

http://developer.venstar.com/
