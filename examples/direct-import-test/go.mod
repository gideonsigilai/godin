module direct-import-test

go 1.24.4

replace github.com/gideonsigilai/godin => ../..

require github.com/gideonsigilai/godin v0.0.0-00010101000000-000000000000

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
