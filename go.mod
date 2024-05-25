module github.com/emicklei/varvoy

go 1.21.5

replace github.com/traefik/yaegi => ../yaegi

replace github.com/traefik-contrib/yaegi-debug-adapter => ../yaegi-debug-adapter

require (
	github.com/emicklei/tre v1.7.0
	github.com/lmittmann/tint v1.0.4
	github.com/traefik-contrib/yaegi-debug-adapter v0.0.0-20220419125716-d76c0ac9f292
	github.com/traefik/yaegi v0.16.1
	golang.org/x/mod v0.17.0
)
