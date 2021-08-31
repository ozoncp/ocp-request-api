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

The compiled binary placed at `bin/ocp-request-api`. To start a local database and other services
run `docker compose up` from repository root. To create all tables run `make migrate`.

### Run tests

To run tests execute `make test` from repository root. It wil run the tests and print coverage report.

### To build and run with Docker

- Build docker image `docker build . -t ocp-request-api`
-
Run `docker run -v <path to config>.yaml:/root/config.yaml  -p 82:82 ocp-request-api /root/ocp-request-api -c config.yaml`

### Config

See below the config example:

```yaml
general:
  write_batch_size: 100 // Controls batch size of multi create endpoint.
db:
  dsn: "dsds" // defines connection to Postresql (in form of Golang's sql DSN).
kafka:
  brokers: localhost:9094  // A comma separate list of Kafka brokers addresses (e.g. host:ip,host:ip)
jaeger:
  agent_host_port: localhost:6831 // Jaeger host and port (e.g. host:ip)

```

The config can be overridden via OCP_REQUEST_<config value path> prefixed env variables. e.g OCP_REQUEST_JAEGER_AGENT_HOST_PORT=localhost:6831 

