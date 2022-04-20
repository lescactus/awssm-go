FROM golang:1.14 as builder

ENV CGO_ENABLED=0

WORKDIR /src

# Copy go.mod and go.sum for cache efficiency
COPY go.* ./

# Download the required modules
RUN go mod download -x

COPY *.go ./

# Build the app
RUN go build -o /out/awssm-go

FROM alpine:3

COPY --from=builder /out/awssm-go /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/awssm-go"]