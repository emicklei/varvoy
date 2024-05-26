```mermaid
sequenceDiagram
    vscode-go ->> varvoy[dap.Session]: protocolmsg (tcp)
    varvoy[dap.Session]-->>varvoy[ProxyAdapter]: Handle
    varvoy[ProxyAdapter]-->>varvoy[ProxySession[dap.Session]]: Forward
    varvoy[ProxySession[dap.Session]]-->>toDebug[dap.Session]: protocolmsg (tcp)
    toDebug[dap.Session]-->>toDebug[dbg.Adapter]: Handle
```
