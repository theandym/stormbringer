# Stormbringer

A simple Heroku-based distributed load testing tool inspired by [Mjölnir](https://github.com/tsykoduk/Mjolnir).


## Quick Start

Clone this repo:

```
git clone https://github.com/theandym/stormbringer
cd stormbringer
```

Create a new Heroku app with the Go buildpack:

```
heroku apps:create [NAME] --buildpack https://github.com/heroku/heroku-buildpack-go
```

Set the following config vars:

```
heroku config:set \
  TARGETS="https://google.com,https://facebook.com,https://twitter.com" \
  WORKERS=8 \
  LENGTH=10000
```

Push the app to Heroku and scale the `stormbringer` process type:

```
git push heroku master
heroku ps:scale stormbringer=1
```

Each dyno will create `$WORKERS` processes. Each process will request the target endpoints `$LENGTH` times, then sleep.

Tail the logs to monitor the progress:

```
heroku logs --tail
```

Scale down the dynos when the test is complete:

```
heroku ps:scale stormbringer=0
```


## Config

```
stormbringer [--curl] [--workers=WORKERS] [--length=LENGTH] TARGETS
```

### General

Whenever possible, config vars should be used to specify the values for arguments and flags. This allows for rapid changes to configuration without modifying the underlying code or configuration file(s).

### Targets

`stormbringer` requires a single argument (`TARGETS`), which should be a comma-delimited list (no spaces) of target endpoints to request.

Example:

```
"https://google.com,https://facebook.com,https://twitter.com"
```

### Options

The following flags provide the ability to modify the default configuration:

  - `--curl`: Switch to `curl` for requests (default: the Go `net/http` package)
  - `--workers`: The number of worker processes to run (default: `8`)
  - `--length` The number of times each worker will request each target endpoint; use `0` to run continuously (default: `10000`)
