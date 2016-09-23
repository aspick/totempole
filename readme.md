# Totempole

windows daemonize helper tool

## how to use

```text
> ./totem [command]

commands:
    install     install Totempole to windows services
    start       start Totempole service
    stop        stop service
    restart     restart service
    reload      reload and apply new config
    status      show status
    ps          show current daemons
    log         show specified daemons stdout
    uninstall   uninstall Totempole from windows services
```


## config format

```yaml
meta: totempole-1.0
daemons:
    - name: 'webapp queue worker'
      ps: 'php artisan queue:worker'
      pwd: 'C:/inetpub/wwwroot/webapp'
      workers: 3

    - cmd: ''
```

- `name` is required
- either `ps` or `cmd ` is required
- `pwd` is option
- `workers` is num of request to launch daemons
