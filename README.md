# rpssl4bu

Task of RPSSL service for Billups.
Main focus is on backend, frontend is just for fun.

# Backend

## About

This server is for playing game of Rock, Paper, Scissors, Lizard, Spock.
The server can also use an external supply of random numbers.

Some functions are not covered with tests for speedup purposes, but main
functionality is tested.

## Running

First change directory to `./backend`.

To build and run the code, first install the dependencies listed in `go.mod`
then build with `go build -o rpssl ./cmd/main.go` and run the created executable.

You can provide arguments:

1. Listen to a different interface/port, for example: — `rpssl --addr 185.34.1.4:342`. Default `:8080`.
2. Use rng provider at http://youraddress.com/provider, — `rpssl --rng http://youraddress.com/provider`. Default — internal provider based on rand.
3. Set log level — `rpssl --log-level debug`. Default — `info`.
4. Set log type to json or text — `rpssl --log-type json`. Default — `text`.

You can combine these parameters as needed.

## Docker run

First change directory to `./backend`.

1. `docker build -t rpssl4bu .`
2. `docker run -e EXTERNAL_RNG=<rng_provider_address> rpssl4bu`

or just `docker run rpssl4bu` to use internal RNG

8080 is exposed

# Frontend

Based on Vue.JS 3.

## Project setup, run

First change directory to `./frontend_vue`.

Use NPM for running.

```
npm install

# dev
npm run serve

# prod
npm run build
```

## Docker run

First change directory to `./frontend_vue`.

1. `docker build -t rpssl4bu_vue .`
2. `docker run rpssl4bu_vue`

80 is exposed

# Both together

There's a docker-compose setup to run both at the same time.
Port 80 is exposed.

# Additional things

## Optional tasks

Global result board and P2P are only working for this server.
