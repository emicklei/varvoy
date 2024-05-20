## varvoy

A Go debugger that allows code modification.

Is intended to be used with the `vscode-go` extension of Microsoft Visual Studio Code.

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

- yaegi tool install

## current limitations

- go mod file cannot have replace