# ===================================================================================
# === Stage 1: Build the Go service code into 'server' exe ==========================
# ===================================================================================
FROM golang:1.14-alpine as go-build

ARG serviceName="SET_ON_COMMAND_LINE"
ARG version="0.0.1"
ARG buildInfo="Local manual builds"

WORKDIR /build

# Install system dependencies
RUN apk update && apk add git gcc musl-dev

# Fetch and cache Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy in Go source files
COPY cmd/$serviceName/ ./service
COPY pkg/ ./pkg

# Now run the build
# Disabling cgo results in a fully static binary that can run without C libs
# Also inject version and build details 
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.version=$version -X 'main.buildInfo=$buildInfo'" \
    -o server ./service

# ================================================================================================
# === Stage 2: Get server exe into a lightweight container =======================================
# ================================================================================================
FROM alpine
WORKDIR /app 

ARG serviceName="SET_ON_COMMAND_LINE"
ARG servicePort=9000

# Copy the Go server binary
COPY --from=go-build /build/server . 

EXPOSE $servicePort
ENV PORT=$servicePort

# This is a trick, we don't really need run.sh
# But some services might have .db files, some don't
COPY cmd/$serviceName/run.sh cmd/$serviceName/*.db ./

# That's it! Just run the server 
CMD [ "./server"]