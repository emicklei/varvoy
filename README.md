## varvoy

Uses yaegi to debug programs from VSCode.

## install

```
go install github.com/emicklei/varvoy/cmd/varvoy@latest
```

## configure

In `settings.json` of the Go VSCode plugin, set an alternative to `dlv`:

```
    "go.alternateTools": {
        "dlv": "/Users/emicklei/go/bin/varvoy"
    }
```

## requirements

- yaegi tool install

## current limitations

- go mod file cannot have replace