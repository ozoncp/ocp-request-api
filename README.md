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

The compiled binary placed at `bin/ocp-request-api`.
To start a local database and other services run `docker compose up` from repository root. To create all tables run `make migrate`.

### Run tests

To run tests execute `make test` from repository root.


### To build and run with Docker

- Build docker image `docker build . -t ocp-request-api`
- Run `docker run -v "OCP_REQUEST_DSN=<pgsql dsn>" -v "OCP_REQUEST_BATCH_SIZE=1000" -v "OCP_KAFKA_BROKERS=kafka:9094" -v "OCP_REQUEST_JAEGER_HOST_PORT=jaeger:6831"  -p 82:82 ocp-request-api`

### ENV variables

- `OCP_REQUEST_DSN` - defines connection to Postresql (in form of Golang's sql DSN). 
- `OCP_REQUEST_BATCH_SIZE` - Controls batch size of multi create endpoint.
- `OCP_KAFKA_BROKERS` - A comma separate list of Kafka brokers addresses (e.g. host:ip,host:ip)
- `OCP_REQUEST_JAEGER_HOST_PORT` - Jaeger host and port (e.g. host:ip)

