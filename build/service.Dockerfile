# ===================================================================================
# === Stage 1: bleep bloop                                  =========================
# ===================================================================================
FROM golang:1.14-alpine as go-build

ARG serviceName="SET_ON_COMMAND_LINE"
ARG version="0.0.1"
ARG buildInfo="Local manual builds"

WORKDIR /build

# Install system dependencies
RUN apk update && apk add git gcc musl-dev

# Fetch and cache Go modules
COPY services/go.mod .
COPY services/go.sum .
RUN go mod download

# Copy in Go source files
COPY services/$serviceName/ ./service
COPY services/common/ ./common

# Now run the build
# Disabling cgo results in a fully static binary that can run without C libs
# Also inject version and build details 
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.version=$version -X 'main.buildInfo=$buildInfo'" \
    -o server ./service

# ================================================================================================
# === Stage 2: Fooo ===================================
# ================================================================================================
FROM alpine
WORKDIR /app 

ARG servicePort=9000

# Copy the Go server binary
COPY --from=go-build /build/server . 

EXPOSE $servicePort
ENV PORT=$servicePort

# That's it! Just run the server 
CMD [ "./server"]