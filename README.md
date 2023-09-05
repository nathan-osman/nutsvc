## nutsvc

This simple Windows service monitors a [NUT](https://networkupstools.org/) server for changes and reacts to power loss events.

### Installation

To install the service, use the `install` argument in an elevated terminal:

    nutsvc.exe install

Then use `start` to run the service.

    nutsvc.exe start
