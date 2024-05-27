## varvoy

**_varvoy is short for "Variable Voyager"_** 

A Go debugger that allows code modification.

Is intended to be used with the `vscode-go` extension of Microsoft Visual Studio Code.

It is build on top of two awesome packages:
- [yaegi](https://github.com/traefik/yaegi)
- [yaegi-debug-adapter](https://github.com/traefik-contrib/yaegi-debug-adapter).

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

## requirements

- Microsoft Visual Studio Code 
- `vscode-go` extension

## current limitations

- no Windows support for now
- project must have a go.mod
- go mod file cannot have replace
