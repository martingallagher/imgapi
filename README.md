# ImgAPI Demo

A simple [CR]UD service for storing and retrieving images using gRPC.

## Instructions

Package management is handled by **Go dep**: <https://github.com/golang/dep>

-   Fetch dependencies using Go dep: `make deps`
-   Use the Makefile to build the CLI: `make build`
-   Use the `-c` flag to load the configuration file: `-c ./config.yml` - this flag can be omitted if `imgapi` is executed within the directory containing `config.yml`
-   Run the server via `./imgapi server`
-   Upload an image via `./imgapi upload testdata/gopher.png` (emits the image ID on success)
-   Download an image via its ID `./imgapi download 0de648a9c8c19264e6cd6a441a867d0989a03929cacec442ad1f0cd192bc9072`
-   Download an image via its ID to a different format `./imgapi download 0de648a9c8c19264e6cd6a441a867d0989a03929cacec442ad1f0cd192bc9072 jpeg`

## Test Scenarios

The demonstration includes unit test and integration test targets:

-   `make test`
-   `make integration_test`
