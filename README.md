# ZeZen
## Public IP to local IP cache that powers your.zencu.be

This HTTP Service accepts requests with your local IP to `/announce` and redirects to the local IP when the same remote IP makes an HTTP request to  `/`.
That allows a Zencube to announce its local IP to this service and `http://your.zencu.be` to redirect to it when request from the same network.

## Prerequisites
Zazen uses [Redis](http://redis.io/) for caching the IP pairs. Make sure you have a Redis server running and available to the Zazen service.
By default a local redis server is assumed but that can be changed using the `--redisaddr` flag.

## Running the service
If you are on a linux x64 host, just run:

```shell
  $ ./zazen
```
which will start the server on port 8080. For a different port (or a specific address) run it with `--addr`:

```shell
  $ ./zazen --addr "127.0.0.1:9000"
```
would start the server listening only on localhost and port 9000.
Other command line options are:

| Flag        | Purpose                             | Default        |
| ----------- |:-----------------------------------:| --------------:|
| --addr      | Address and Port for HTTP server    | :8080          |
| --redisaddr | Address of the redis server         | 127.0.0.1:6379 |
| --ttl       | Time to live for the cach (seconds) | 3600           |
| --db        | The Redis database to use           | 0              |

## Building from source
With Go 1.1+ available run:

```shell
  $ make build
```

## Running it for development
```shell
  $ make run
```

## License
This software is MIT licensed. See LICENSE for details.
