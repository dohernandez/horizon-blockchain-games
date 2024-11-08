# Contributing

* If you are a new contributor see: [Steps to Contribute](#steps-to-contribute).

* Relevant coding style guidelines are the [Go Code Review
  Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)
  and the _Formatting and style_ section of Peter Bourgon's [Go: Best
  Practices for Production
  Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style).


## Steps to Contribute

Before you start contributing, make sure you have the following tools installed:

* [Go](https://golang.org/dl/) version 1.23.x or greater.
* [Docker](https://docs.docker.com/get-docker/) version 20.10.x or greater.
* [Docker Compose](https://docs.docker.com/compose/install/) version 1.29.x or greater.
* [golangci-lint](https://github.com/golangci/golangci-lint/releases) version v1.61.x or greater.
* [gci](https://github.com/daixiang0/gci/releases) version v0.13.x or greater.
* [gofumpt](https://github.com/mvdan/gofumpt/releases) version v0.7.x or greater.
* [mockery](https://github.com/vektra/mockery/releases) version v2.46.x or greater.

For quickly installing the tools (`golangci-lint`, `gci`, `gofumpt`, `mockery`) do:

```bash
make tools
```

For quickly compiling and testing your change(s) do:

```bash
make test         # Make sure all the tests pass before you commit and push :)
```

For linting the code do:

```bash
make lint        # Make sure your change(s) follow our coding standards.
```

For checking the commit do:

```bash
make check        # Make sure test and lint pass before you commit and push :)
```

For more tools and options see the [Makefile](Makefile).

```bash
make

Usage

  tools:                Install all require tools to work with the project
  test:                 Run tests
  check:                Check test and linter the change.
  install-linter:       Check/install golangci-lint tool
  install-gci:          Check/install gci tool
  install-gofumpt:      Check/install gofumpt tool
  lint:                 Check with golangci-lint
  fix-lint:             Apply goimports and gofmt
  test-unit:            Run unit tests
  test-unit-multi:      Run unit tests multiple times, use `UNIT_TEST_COUNT=10 make test-unit-multi` to control count
  install-mockery:      Check/install mockery tool
  build-linux:          Build Linux binary
  build-darwin-amd:     Build macOS intel binary
  build-darwin-arm:     Build macOS Apple M1 binary
  build:                Build binary
  build-image:          Build docker image
  release-assets:       Build and compress binaries for release assets.
```


## Pull Request

* Branch from the main branch and, if needed, rebase to the current main branch before submitting your pull request. If it doesn't merge cleanly with main you may be asked to rebase your changes.

* Commits should be as small as possible, while ensuring that each commit is correct independently (i.e., each commit should compile and pass tests).

* Add tests relevant to the fixed bug or new feature.

## Dependency management

Avoid introducing external dependencies without a good reason, but if so the project uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies on external packages. This requires a working Go environment with version 1.12 or greater installed (version of the project 1.23.3).

All dependencies are vendored in the `vendor/` directory.

To add or update a new dependency, use the `go get` command:

```bash
# Pick the latest tagged release.
go get example.com/some/module/pkg

# Pick a specific version.
go get example.com/some/module/pkg@vX.Y.Z
```

Tidy up the `go.mod` and `go.sum` files and copy the new/updated dependency to the `vendor/` directory:


```bash
# The GO111MODULE variable can be omitted when the code isn't located in GOPATH.
GO111MODULE=on go mod tidy

GO111MODULE=on go mod vendor
```

You have to commit the changes to `go.mod` and `go.sum` before submitting the pull request.


Happy coding!!!