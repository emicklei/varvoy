module github.com/emicklei/varvoy

go 1.22

toolchain go1.22.4

replace github.com/traefik/yaegi => github.com/emicklei/yaegi v0.2.1

// replace github.com/traefik/yaegi => ../yaegi

// replace github.com/traefik-contrib/yaegi-debug-adapter => ../yaegi-debug-adapter

// replace github.com/emicklei/structexplorer => ../structexplorer

require (
	github.com/emicklei/structexplorer v0.0.0-20240704124131-4ccb1c7cab9c
	github.com/traefik-contrib/yaegi-debug-adapter v0.0.0-20240606200100-1922144a1da7
	github.com/traefik/yaegi v0.16.1
	golang.org/x/mod v0.17.0
)

require github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
