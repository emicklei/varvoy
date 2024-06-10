module github.com/emicklei/varvoy

go 1.21.5

replace github.com/traefik/yaegi => github.com/emicklei/yaegi v0.2.1

// replace github.com/traefik/yaegi => ../yaegi

// replace github.com/traefik-contrib/yaegi-debug-adapter => ../yaegi-debug-adapter

require (
	github.com/traefik-contrib/yaegi-debug-adapter v0.0.0-20240606200100-1922144a1da7
	github.com/traefik/yaegi v0.16.1
	golang.org/x/mod v0.17.0
)
