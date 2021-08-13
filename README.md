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
The compiled binary placed at `bin/ocp-request-api`


### To build and run with Docker

- Build docker image `docker build . -t ocp-request-api`
- Run `docker run -p 82:82 ocp-request-api`
