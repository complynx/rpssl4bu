# rpssl4bu

Task of RPSSL service for Billups.
Main focus is on backend, frontend is just for fun.
Thus only main backend task was properly structured and covered with tests.
Frontend and some optional tasks are not covered with tests.

## Optional tasks

Optional tasks were done as well:
1. GUI uses Vue 3 (see below)
2. Scoreboards are implemented for local and global scores
3. Every scoreboard may be reset
4. P2P is implemented using WebSockets and with separate scoreboard
5. Multiple users can easily and separately play on the server
6. Dockerfile-s are provided, as well as full docker-compose build
7. Some fun things were added with the help of new technologies

# Backend

## About

This server is for playing game of Rock, Paper, Scissors, Lizard, Spock.
The server can also use an external supply of random numbers (though this
might slow down some functions). To use external provider, please provide
necessary arguments or variables. 

Some features are not covered with tests for speedup purposes, but main
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

Port 8080 is exposed by default

# Frontend

Based on Vue.JS 3.
This frontend is created in one component because frontend isn't
the main part of the task. No testing either.
Some features will stop working for different backend servers.

## Project preparation

If you want to run together with backend in docker-compose, nothing
needs to be done. If you want to run against different backend server
Fix value of `backendServer` property in `./frontend_vue/src/App.vue`.
You can also change backend server on the go later from GUI.

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

Don't forget to provide proper backend server (from GUI or by fixing App.vue).

## Docker run

First change directory to `./frontend_vue`.

1. `docker build -t rpssl4bu_vue .`
2. `docker run rpssl4bu_vue`

There are two variable that you could use:
1. FRONTEND_LOCATION — final location of the ./ for vue to build from
2. BACKEND_URL — location of the backend (sets up default `backendServer`)

Port 80 is exposed by default.

# Both together

There's a docker-compose setup to run both backend and frontend at the same time.
To run both together, just run `docker-compose build` and `docker-compose up`.
Frontend will be mounted in `/` and backend will be mounted in `/backend`.

Port 80 is exposed.
