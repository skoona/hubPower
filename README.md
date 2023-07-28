# HubPower

A application to display voltage and wattage from a Hubitat Home Automatin Controller.




### Project Layout

    Enable debug logging via environment var: `export SKN_DEBUG="true"`
    Specify Hubitat Hub's IpAddress via environment var: `export HUBITAT_IP="10.100.1.41"`
    Specify this machine IP via environment var: `export TRUSTED_IP="10.100.1.5"`


```text
├── FyneApp.toml
├── LICENSE
├── Makefile
├── README.md
├── bin
│   └── hubPower
├── commons
│   ├── common.go
│   ├── imageResources.go
│   └── images.go
├── docs
│   ├── allDevices.json
│   ├── allDevicesFull.json
│   ├── deviceCapabilites7.json
│   ├── deviceCapabilities3.json
│   └── deviceHistory7.json
├── entities
│   ├── graphaverage.go
│   ├── hubdevices.go
│   └── hubhost.go
├── go.mod
├── go.sum
├── interfaces
│   ├── configuration.go
│   ├── graphpointsmoothing.go
│   ├── hubprovider.go
│   ├── provider.go
│   ├── service.go
│   └── viewprovider.go
├── main.go
├── providers
│   ├── config.go
│   └── hubitatprovider.go
├── resources
│   ├── apcupsd.png
│   └── gapc_prefs.png
├── services
│   └── service.go
└── ui
    ├── detailed.go.notReady
    ├── menus.go
    ├── monitor.go
    ├── overview.go
    ├── preferences.go
    └── viewprovider.go
```

### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request


### LICENSE
The application is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
