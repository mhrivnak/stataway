# Stataway

Stataway is a service that runs in your home, determines if anyone is in your
home, and sets your thermostat to its "home" or "away" mode accordingly.

The service enables laziness for those who forget to manually set their thermostat
to "home" or "away", but it does not get in the way of doing so. If you are
paying the bill at a restaurant and use a mobile app to set your thermostat to
"home", this service will not "fight" your state change.

Put another way, as with all good home automation, this service endeavors to
be helpful without being a dependency or getting in the way.

## Settings

In addition to the detector-specific and thermostat-specific settings below,
these configuration settings may be set as environment variables:

`HOME_LATITUDE`: latitude of your home location in decimal format

`HOME_LONGITUDE`: longitude of your home location in decimal format

## Thermostats

The only supported thermostats currently are those made by Venstar and that
support its [local API](http://developer.venstar.com/).

Support for additional thermostats is desired. Please be in touch if you would
like to assist in adding one.

### Settings

`VENSTAR_URL`: URL to the thermostat's API

## Detectors

Whether anyone is in your home is determined by "detectors". Currently there
is only one detector, but the system is designed for multiple.

A detector sends a "trigger" notification when it believes the state should
change from "home" to "away" or vice versa. Those notifications get received by
a control loop that determines if and how to make a change to a thermostat's
state.

### Google Location Sharing

This detector uses Google location sharing to determine if someone is close
enough to your home for the thermostat to be in "home" mode. It polls for
the locations every 30 seconds.

#### Inner and Outer Shells

The detector uses two distances to determine when to send a trigger. When
changing from a state where any shared location is inside the outer shell to a
state where all shared locations are outside the outer shell, a trigger is sent
for the "away" state.

After this detector has triggered "away", it waits for at least one location to
appear within the inner shell before triggering "home".

Having the inner and outer shells helps avoid thrashing between the two states
in a case where someone was near the boundary in a single-shell design.

What distances you choose will depend on your normal travel habits, and when
you want your thermostat to change states. Early testing has used .7 and .5 km. 
Those distances are good enough to cause the thermostat to go to "away" most of
the time when the house is unoccupied for an hour or more. But the distances are small
enough that the service rarely notices that occupants are close to home until they
are actually at home, so it does not give the HVAC an opportunity to pre-condition
the space.

In practice, it seems that when a device first connects to the wifi network at
home, that causes the device to report a location change. That is usually good
enough for your system to engage heating or cooling before you make it inside,
but only by a minute or less.

#### Getting Started

1. Create a new Google account only for this purpose.
1. Have every person whose location you want to consider share their location with
   that new account through Google Maps.
1. Set the below environment variables when running stataway.

`GOOGLE_USERNAME`: The full username for a Google account, usually an email address.

`GOOGLE_PASSWORD`: The password for the account.

`GOOGLE_INNER_KM`: The inner shell distance in kilometers, decimal values supported.

`GOOGLE_OUTER_KM`: The outer shell distance in kilometers, decimal values supported.

### Future Detectors

No, we haven't figured out how to detect the future. These are ideas for detectors that
could be implemented in the future.

#### Owntracks

[Owntracks](http://owntracks.org/) could be used as a location source similar
to Google. It has privacy advantages. It also is meant to be integrated with,
as opposed to Google location sharing data that we have to scrape from an
undocumented API.

A downside is that you may need to provide a service that listens on the
Internet for devices to report location changes.

#### LAN

A detector could look for the presence of devices on the local network.

Depending on how this gets implemented, it can be a bad idea to ping mobile
devices repeatedly and regularly, because it can prevent them from staying in a
power-saving mode. If polling for their presence, it should only be done while
the current known state is "away", and should be used to detect the arrival of
a person at home. Other detectors could then be used to determine when to
change from "home" to "away".
