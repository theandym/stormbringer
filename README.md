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

Each dyno will create `$WORKERS` processes. Each process will cURL the target endpoints `$LENGTH` times, then die.

Tail the logs to monitor the progress:

```
heroku logs --tail
```

To rerun the test, scale down and up or restart the dynos.


## Config Vars

  - `TARGETS`: A comma delimited list (no spaces) of target endpoints to cURL (ex. `"https://google.com,https://facebook.com,https://twitter.com"`)
  - `WORKERS`: The number of worker processes to run (default: `8`)
  - `LENGTH` The number of times each worker will cURL each endpoint; use `0` to run continuously (default: `10000`)
