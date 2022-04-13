# HostMonitor

A simple utility to monitor host availability by given config.

### Usage:

```bash
hostmonitor watch [flags]
```

### Flags:

-c, --configUrl string Url of config containing url list (default "https://gist.githubusercontent.com/ddosukraine2022/f739250dba308a7a2215617b17114be9/raw/mhdos_targets_tcp.txt")

-h, --help help for watch

-t, --requestTimeout int Request timeout (default 5)

-i, --testInterval int Interval in seconds between test updates (default 10)
