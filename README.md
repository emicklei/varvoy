## varvoy

A Go debugger that allows code modification.

It is build on top of two awesome packages [yaegi](https://github.com/traefik/yaegi) and [yaegi-debug-adapater](https://github.com/traefik-contrib/yaegi-debug-adapter).

## install

```
go install github.com/emicklei/varvoy/cmd/varvoy@latest
```

## configure

Is intended to be used with the `vscode-go` extension of Microsoft Visual Studio Code.
In `settings.json` of the Go VSCode plugin, set an alternative to `dlv`.
Use `which varvoy` to find the absolute path.

```
    "go.alternateTools": {
        "dlv": "/Users/emicklei/go/bin/varvoy"
    }
```

## requirements

- Microsoft Visual Studio Code 
- `vscode-go` extension

## current limitations

- linux only for now
- go mod file cannot have replace
- project must have a go.mod