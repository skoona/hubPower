# ggApcMon

A go refactor of gapcmon from source forge, originally written in C/Gtk2.




### Project Layout

Enable debug logging via environment var: `export GAPC_DEBUG="true"`


```text
├── LICENSE
├── README.md
├── bin
│   └── ggapcmon
├── cmd
│   ├── cli
│   │   └── main.go
│   └── gui
│       └── main.go
├── docs
├── go.mod
└── internal
    ├── entities
    │   └── apchosts.go
    ├── interfaces
    │   ├── apcprovider.go
    │   ├── provider.go
    │   ├── service.go
    │   └── viewprovider.go
    ├── providers
    │   ├── apcprovider.go
    │   └── viewprovider.go
    ├── services
    │   └── service.go
    └── ui
        ├── monitor.go
        └── settings.go
```

### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request


### LICENSE
The application is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
