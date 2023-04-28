# rpssl4bu

Task of [RPSSL](https://bigbangtheory.fandom.com/wiki/Rock,_Paper,_Scissors,_Lizard,_Spock) service for Billups.
Main focus is on backend written in Go, frontend is just for fun.
Thus only main backend task was properly structured and covered with tests.
Frontend and some optional backend tasks are not covered with tests and
may present bad development practices.

[Check it out live](https://complynx.net/rpssl/)

## Optional tasks

Optional tasks were done as well:
1. GUI uses Vue 3 (see below for details)
2. Scoreboards are implemented for local and global scores
3. Every scoreboard may be reset
4. P2P is implemented using WebSockets and with a separate scoreboard per game
5. Multiple users can easily and separately play on the server
6. Dockerfile-s are provided, as well as full docker-compose build
7. Some fun things were added with the help of new technologies (see [live version](https://complynx.net/rpssl/))

Main task was:

# Backend

## About

This server is for playing game of Rock, Paper, Scissors, Lizard, Spock.
The server can  use an external supply of random numbers (though this
might slow down some functions). To use external provider, please provide
necessary arguments to the executable. Alternatively you can provide them
through environment variable for Docker.

Some features are not covered with tests for speedup purposes, but main
task functionality — three required endpoints and use of external RNG — are tested.

## Pure setup and run

First change directory to `./backend`.

To build and run the code, first install the dependencies listed in `go.mod`
then build with `go build -o rpssl ./cmd/main.go` and run the created executable.

You can provide arguments:

1. Listen to a different interface:port, for example: — `rpssl --addr 185.34.1.4:342`. Default is `:8080`.
2. Use external rng provider at http://youraddress.com/provider, — `rpssl --rng http://youraddress.com/provider`. Default — internal provider based on library `rand`.
3. Set log level — `rpssl --log-level debug`. Default — `info`.
4. Set log type to json or text — `rpssl --log-type json`. Default — `text`.

You can combine these parameters as needed.

Alternatively you can:

## Docker run

First change directory to `./backend`.

1. `docker build -t rpssl4bu .`
2. `docker run -e EXTERNAL_RNG=<rng_provider_address> rpssl4bu`

or just `docker run rpssl4bu` to use internal RNG

Port 8080 is exposed by default.

To experience the full powers of the backend, you need to use specially crafted...

# Frontend

This frontend, being based on Vue.JS 3, is created as solo component because
the frontend isn't the main part of the task. No testing was added either.

The frontend can be bound to different backend servers, though some features
will stop working for non-authentic backend servers.

## Frontend preparation

If you want to run it together with the backend in docker-compose, nothing
needs to be done — just use the provided `docker-compose.yml` file.

If you want to run against a different backend server, you can:
1. Fix the value of `backendServer` property in `./frontend_vue/src/App.vue`
2. Change backend server on the go directly from the GUI
3. Set up `VUE_APP_BACKEND_URL` in `.env` before building Vue application
3. Set up the `BACKEND_URL` environment variable when building a Docker container

## Pure setup and run

First change directory to `./frontend_vue`.

Ensure you have the latest Node.JS and NPM installed, then:

```
echo VUE_APP_BACKEND_URL=http://your.backend/server/>.env
npm install
npm run serve
```
or alternatively
```
npm run build
```

As another way to run, you can use

## Docker run

First change directory to `./frontend_vue`.

1. `docker build -t rpssl4bu_vue .`
2. `docker run rpssl4bu_vue`

There are two environment variables that you could set up for the build of the container:
1. `FRONTEND_LOCATION` — the prefix at which the Vue application will be located, must be set if your GUI will not be located in the root server directory
2. `BACKEND_URL` — address of the backend service endpoints root

Port 80 is exposed by default.

But probably you'd like to run

# Both together

There's a `docker-compose.yml` provided to run both backend and frontend at
the same time, already bundled up. To run both together, just run
`docker-compose build` and `docker-compose up`.
Frontend will be mounted in `/` and backend will be mounted in `/backend/`
of your server.

Under the hood it will use a couple of NGINX servers, the primary one is set up
in `./configs/nginx/conf.d/main.conf`.

Port 80 is exposed.
