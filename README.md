# Issuer ID middleware

The `issuerid` package provides a middleware for HTTP servers in Go that generates a unique identifier (UUID) based on the client's real IP address. This identifier can be used for tracking or logging purposes in applications where identifying the source of requests is important.

## Features
 * Extracts the real IP address of the client from HTTP headers.
 * Generates a UUID based on the client's IP address using SHA-1 hashing.
 * Provides a context-based way to retrieve the generated UUID throughout the request lifecycle.

## Installation

To use this library in your Go project, you can import it using the following command:

```bash
go get github.com/junkofuruto/issuer-id
```

## Usage

### Middleware

To use the `issuerid` middleware, wrap your HTTP handler with it. The middleware will extract the client's real IP address and generate a UUID, which will be stored in the request context.

```go 
package main

import (
	"net/http"
	"github.com/junkofuruto/issuerid"
)

func main() {
	http.Handle("/", issuerid.IssuerId(http.HandlerFunc(yourHandler)))
	http.ListenAndServe(":8080", nil)
}

func yourHandler(w http.ResponseWriter, r *http.Request) {
	issuerID := issuerid.GetIssuerID(r.Context())
	w.Write([]byte("Your Issuer ID: " + issuerID))
}
```

### Retrieving the Issuer ID

```go
func yourHandler(w http.ResponseWriter, r *http.Request) {
	issuerID := issuerid.GetIssuerID(r.Context())
	// Use the issuerID for logging or tracking
}
```

## How It Works
1. The middleware checks for the client's real IP address using the following headers:
   - `True-Client-IP`
   - `X-Forwarded-For`
   - `X-Real-IP`
2. If a valid IP address is found, it generates a UUID from the IP address using SHA-1 hashing.
3. The generated UUID is stored in the request context, allowing it to be accessed later in the request handling process.
4. If no valid IP address is found, a zero UUID (00000000-0000-0000-0000-000000000000) is used.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue for any enhancements or bug fixes.

## Acknowledgments
This package utilizes the github.com/google/uuid package for UUID generation.
Special thanks to the Go community for their contributions and support.