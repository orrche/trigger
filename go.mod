module github.com/orrche/trigger

go 1.14

require (
	github.com/gofrs/uuid v4.1.0+incompatible // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/pat v1.0.1
	github.com/ivanbeldad/kasa-go v0.0.0-20201031100518-9b33fa73f8a7
	github.com/orrche/trigger/triggerkasa v0.0.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/orrche/trigger/triggerkasa v0.0.0 => ./triggerkasa
