## varvoy

**_varvoy is short for "Variable Voyager"_** 

A Go debugger.
 
It is build on top of two awesome packages:
- [yaegi](https://github.com/traefik/yaegi)
- [yaegi-debug-adapter](https://github.com/traefik-contrib/yaegi-debug-adapter).

## features

- works with the `vscode-go` extension of Microsoft Visual Studio Code
- imports generated stubs for required Go modules

## goal

- handle code modification during debugging session
- evaluate expression in a debugger session
- restart from stack frame

## install

```
go install github.com/emicklei/varvoy/cmd/varvoy@latest
```

## configure

In `settings.json` of the Go VSCode plugin, set an alternative to `dlv`.
Use `which varvoy` to find the absolute path.

```
    "go.alternateTools": {
        "dlv": "/Users/emicklei/go/bin/varvoy"
    }
```

## design

See [Desgin](./doc/README.md)


## current limitations

- no Windows support for now
- project must have a go.mod
- go mod file cannot have replace
