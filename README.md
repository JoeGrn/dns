# Simple DNS Server

A lightweight DNS server implementation written in Go.

## Overview

This DNS server provides basic DNS functionality, including:

- UDP-based DNS message handling
- Support for standard DNS queries
- Lightweight and containerized deployment

## Project Structure

```
├── Dockerfile       # Docker configuration for containerization
├── go.mod           # Go module dependencies
├── Makefile         # Build automation
├── cmd/             # Command-line entrypoints
│   └── dns/         # Main DNS server binary
│       └── main.go  # Entry point for the DNS server
└── pkg/             # Core DNS implementation
    ├── answer.go    # DNS answer section handling
    ├── dns.go       # Core DNS functionality
    ├── flags.go     # DNS flag handling
    ├── header.go    # DNS header implementation
    ├── question.go  # DNS question section handling
    └── server.go    # Server implementation
```

## Getting Started

### Prerequisites

- Go 1.18+ (or Docker for containerized use)
- Make (optional, for using the provided Makefile)

### Building from Source

Clone the repository and build the server:

```bash
# Build the binary
make build

# Run tests
make test

# Format the codebase
make fmt
```

### Running the Server

```bash
# Using make
make run

# Or run the binary directly
./dist/dns
```

The server will start listening for DNS queries on localhost (127.0.0.1) port 2053 by default.

## Docker Support

### Building the Docker Image

```bash
make docker-build
# or
docker build -t dns .
```

### Running with Docker

```bash
make docker-run
# or
docker run -p 2053:2053/udp dns
```

## Testing the Server

You can test your DNS server using tools like `dig`:

```bash
dig @127.0.0.1 -p 2053 example.com
```

## Development

### Making Changes

1. Modify the code in the `pkg/` directory for DNS functionality
2. Use `make test` to verify your changes work as expected
3. Use `make fmt` to format the code according to Go standards

## License

[MIT License](LICENSE) - Feel free to use and modify this code for your own projects.

## Acknowledgements

This project is created as a learning exercise to understand DNS protocols and Go network programming.
