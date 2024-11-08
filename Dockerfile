# syntax=docker/dockerfile:1.2

# --- BEGINING OF BUILDER

FROM golang:1.23.3-bookworm AS builder

#ARG GH_ACCESS_TOKEN

WORKDIR /go/src/github.com/dohernandez/sequence

# This is to cache the Go modules in their own Docker layer by
# using `go mod download`, so that next steps in the Docker build process
# won't need to download modules again if no modules have been updated.
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . ./

# Build rpc binary and cli binary
RUN make build

# --- END OF BUILDER

FROM debian:bookworm

RUN groupadd -r sequence && useradd --no-log-init -r -g sequence sequence
USER sequence

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=sequence:sequence /go/src/github.com/dohernandez/sequence/bin/sequencecli /bin/sequence

ENTRYPOINT ["sequence"]
