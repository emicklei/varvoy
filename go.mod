module github.com/emicklei/varvoy

go 1.21.5

replace github.com/traefik/yaegi => github.com/emicklei/yaegi v0.1.0

// replace github.com/traefik/yaegi => ../yaegi

replace github.com/traefik-contrib/yaegi-debug-adapter => github.com/emicklei/yaegi-debug-adapter v0.1.0

require (
	github.com/traefik-contrib/yaegi-debug-adapter v0.1.0
	github.com/traefik/yaegi v0.16.1
	golang.org/x/mod v0.17.0
)
