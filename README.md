# ocp-request-api

Student requests API. Currently supports:

- Create new request
- Return detailed request information
- Remove request
- List requests

The service accepts gRPC connections at port 82 and HTTP at 8082.

### To build locally

- Install `protoc`. See instruction [here](https://grpc.io/docs/protoc-installation/)
- Build:

```shell
git clone https://github.com/ozoncp/ocp-request-api.git
cd ocp-request-api
make build
```


The compiled binary placed at `bin/ocp-request-api`. The app requires a database DSN to be passed as `OCP_REQUEST_DSN` env variable:

```shell
OCP_REQUEST_API="postgres://ocp-request:12345@localhost:5432/ocp-request?sslmode=disable" bin/ocp-request-api 
```

To start a local database run `docker compose up` from repository root. It starts dev Postgresql server locally.
To create all tables run `make migrate`.

### Run tests

To run tests execute `make test` from repository root.


### To build and run with Docker

- Build docker image `docker build . -t ocp-request-api`
- Run `docker run -p 82:82 ocp-request-api`

