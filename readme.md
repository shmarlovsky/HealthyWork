Small app which sits in tray and send reminders to work sitting or standing (if you have adjustable desk).
Intervals and start state are configured in `conf.yaml` file.

To avoid opening a console at application startup, use these compile flags:  
`go build -ldflags -H=windowsgui`