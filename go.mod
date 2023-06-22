module github.com/cyverse-de/search

go 1.16

replace github.com/cyverse-de/querydsl => ../querydsl

require (
	github.com/cyverse-de/configurate v0.0.0-20171005230251-9b512d37328e
	github.com/cyverse-de/go-mod/otelutils v0.0.2
	github.com/cyverse-de/querydsl v0.0.0-20190124215511-d0881ab0f52c
	github.com/fsnotify/fsnotify v1.4.3-0.20170329110642-4da3e2cfbabc // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/hcl v0.0.0-20171017181929-23c074d0eceb // indirect
	github.com/magiconair/properties v1.7.5-0.20171031211101-49d762b9817b // indirect
	github.com/mitchellh/mapstructure v1.4.2
	github.com/olivere/elastic/v7 v7.0.12
	github.com/pelletier/go-toml v1.0.2-0.20171024211038-4e9e0ee19b60 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.0.4-0.20171118124223-95cd2b9c79aa
	github.com/spf13/afero v1.0.0 // indirect
	github.com/spf13/cast v1.1.0 // indirect
	github.com/spf13/jwalterweatherman v0.0.0-20170901151539-12bd96e66386 // indirect
	github.com/spf13/pflag v1.0.1-0.20171106142849-4c012f6dcd95 // indirect
	github.com/spf13/viper v1.0.1-0.20171109205716-4dddf7c62e16
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.31.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.31.0
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
)
