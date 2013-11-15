# Rehearse
Rehearse provides a simple way to develop single page apps, completely decoupled from any back end. Rehearse will serve your static assets and stub server responses.

## Usage
To use Rehearse, simply download the binary and start the server by running `rehearse` with any of the options given below:

	-address -- The address for the stub server to bind to.
	-port -- The port for the stub server to bind to (defaults to 3333).
	-path -- The path to the static assets you would like Rehearse to serve.

To stub out an endpoint, simply `POST` to `/stubs` with the following JSON:

	{
		"method": "GET", // HTTP method
		"path": "/foo", // Path you wish to stub
		"body": "{\"foo\": \"bar\" }" // Response body of stubbed endpoint
	}

## Contributing
Fork this repo and run tests with:

	go test

Create a feature branch, write your tests and code and submit a pull request.

## License
MIT