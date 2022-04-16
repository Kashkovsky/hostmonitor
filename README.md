# HostMonitor

A simple utility to monitor host availability by given config.

### Usage:

#### Run as a CLI tool to display test metrics in terminal:

```bash
hostmonitor watch [flags]
```

#### Run as a web server:

Web UI will be available at `http://localhost:8080`.
The default port can be changed with a `-p` flag.

```bash
hostmonitor serve [flags]
```

### Flags:

**-c** or **--configUrl**: [string] Url of config containing url list (default "https://raw.githubusercontent.com/Kashkovsky/hostmonitor/main/itarmy_targets.txt")

**-h** or **--help**: help for a command (e.g. `hostmonitor serve -h`)

**-t** or **--requestTimeout**: [int] Request timeout in seconds (default 5)

**-i** or **--testInterval**: [int] Interval between test updates in seconds (default 10)

**-u** or **--updateInterval**: [int] Config update interval in seconds (default 600)

**-p** or **--port**: [int] Server port (default 8080, **for web server only**)
