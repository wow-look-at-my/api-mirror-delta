# api-mirror-delta

A development tool for [api-mirror](https://github.com/wow-look-at-my/api-mirror). It sits in front of the real API and one or more mirrors, and sends each request to all of them. It then reports how every mirror answer differs from the real one.

## Use

Proxy mode — point a client at the delta instead of the mirror:

```sh
api-mirror-delta serve -c delta.yaml
GITHUB_API_URL=http://127.0.0.1:8090 gh api repos/cli/cli
curl -s 'http://127.0.0.1:8090/_delta/summary?format=text'
```

One-shot mode — compare a list of paths and exit 1 on any difference:

```sh
api-mirror-delta compare --truth https://api.github.com --mirror local=http://127.0.0.1:8080 \
	-H "Authorization: token $GITHUB_TOKEN" --ignore '$..*_url' /repos/cli/cli /users/octocat
```

## What it compares

- The status code, the headers you list, and the body. A JSON body is diffed by structure. Anything else byte by byte.
- Each target's own base URL is replaced with `{base}` first, so self-links agree.
- Only `GET` and `HEAD` fan out by default. Other methods go to the served target alone, so a write runs once.

## Ignoring differences

Each rule selects differences by JSONPath (`$..url`, `$.items[*].id`, `$..*_url`), by header, by kind (`missing`, `extra`, `changed`, `type`, `status`, `header`, `body`, `error`), and optionally by route and mirror. The summary counts what each rule hides and flags a rule that matched nothing.

[`samples/github.yaml`](samples/github.yaml) is a starting config. [`docs/config.md`](docs/config.md) is the full reference.

## Build

`go-toolchain` builds, tests, and runs the `dats/` suite.
