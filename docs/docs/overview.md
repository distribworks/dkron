# REST API

You can communicate with dkron using a RESTful JSON API over HTTP. dkron nodes usually listen on port `8080` for API requests. All examples in this section assume that you've found a running leader at `dkron-node:8080`.

dkron implements a RESTful JSON API over HTTP to communicate with software clients. dkron listens in port `8080` by default. All examples in this section assume that you're using the default port.

Default API responses are unformatted JSON add the `pretty=true` param to format the response.
