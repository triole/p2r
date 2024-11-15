# P2R ![build](https://github.com/triole/p2r/actions/workflows/build.yaml/badge.svg) ![test](https://github.com/triole/p2r/actions/workflows/test.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Config example](#config-example)

<!-- /toc -->

## Synopsis

A very simple and basic tool helping to synchronise and push files to remote machines utilising command line tools. These commands, the local and the remote target are configured inside a yaml file. If executed p2r looks for a config file in the current directory that should be named `p2r.yaml`. Other config file can be specified using the `-f` flag.

## Config example

```go mdox-exec="tail -n +2 examples/p2r.yaml"
sync_steps:
  - cmd: ["rsync", "-av", "--delete", "--chown={{.user}}:admins", "--exclude=acme"]
    local: {{.confdir}}/
    remote: remote_machine:/etc/whatever/
  - cmd: ["scp"]
    local: {{.confdir}}/
    remote: remote_machine:/etc/whatever/
```

## Help

```go mdox-exec="sh/display_help.sh"

push or pull remote file systems

Arguments:
  [<action>]    action to perform, can be: [info,push,pull,list]

Flags:
  -h, --help                      Show context-sensitive help.
  -f, --config="$(pwd)/p2r.yaml"
                                  config file
  -n, --dry-run                   only print commands what would have been
                                  executed
  -m, --rsync-dry-run             enable rsync dry runs
      --log-file="/dev/stdout"    log file
      --log-level="info"          log level, can be: [trace,debug,info,error]
      --log-no-colors             disable output colours, print plain text
      --log-json                  enable json log, instead of text one
  -V, --version-flag              display version
```
