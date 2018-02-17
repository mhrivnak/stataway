This is basically a port of this python library: https://github.com/costastf/locationsharinglib

## Usage

Start by creating a dedicated google account and sharing one or more device
locations with it through the "Location Sharing" feature of Google Maps. Using
credentials to access that account, this library will retrieve the name,
latitude and longitude of each person whose location is shared.

It cannot retrieve the location of the account it logs in as.
