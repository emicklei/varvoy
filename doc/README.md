
![flow](varvoy_flow.png "Varvoy desgin")

```mermaid
sequenceDiagram
    vscode-go ->> varvoy[dap.Session]: protocolmsg (tcp)
    varvoy[dap.Session]-->>varvoy[ProxyAdapter]: Handle
    varvoy[ProxyAdapter]-->>varvoy[ProxySession]: Forward
    varvoy[ProxySession]-->>_debug_bin[dap.Session]: protocolmsg (tcp)
    _debug_bin[dap.Session]-->>_debug_bin[dbg.Adapter]: Handle
```
