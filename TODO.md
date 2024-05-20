
```
// Start a process:
cmd := exec.Command("sleep", "5")
if err := cmd.Start(); err != nil {
    log.Fatal(err)
}

// Kill it:
if err := cmd.Process.Kill(); err != nil {
    log.Fatal("failed to kill process: ", err)
}
```

ps aux | grep varvoy | awk '{print $2}' | xargs kill