# ===================================================================================
# === Stage 1: Build the Go service code into 'server' exe ==========================
# ===================================================================================
FROM golang:1.14-alpine as go-build

ARG SERVICE_NAME="SET_ON_COMMAND_LINE"
ARG VERSION="0.0.1"
ARG BUILD_INFO="Not provided"
ARG CGO_ENABLED=0

WORKDIR /build

# Install system dependencies, if CGO_ENABLED=1
RUN if [[ $CGO_ENABLED -eq 1 ]]; then apk update && apk add gcc musl-dev; fi

# Fetch and cache Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy in Go source files
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

# Now run the build
# Inject version and build details, to be available at runtime 
RUN GO111MODULE=on CGO_ENABLED=$CGO_ENABLED GOOS=linux \
go build \
-ldflags "-X main.version=$VERSION -X 'main.buildInfo=$BUILD_INFO'" \
-o server github.com/benc-uk/dapr-store/cmd/${SERVICE_NAME}

# ================================================================================================
# === Stage 2: Get server exe into a lightweight container =======================================
# ================================================================================================
FROM alpine
WORKDIR /app 

ARG SERVICE_NAME="SET_ON_COMMAND_LINE"
ARG SERVICE_PORT=9000

# Copy the Go server binary
COPY --from=go-build /build/server . 

EXPOSE $SERVICE_PORT
ENV PORT=$SERVICE_PORT

# This is a trick, we don't really need run.sh
# But some services might have .db files, some don't
COPY cmd/$SERVICE_NAME/run.sh cmd/$SERVICE_NAME/*.db ./

# That's it! Just run the server 
CMD [ "./server"]