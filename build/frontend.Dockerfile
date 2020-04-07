# ===================================================================================
# === Stage 1: Build static serving host ============================================
# ===================================================================================
FROM golang:1.14-alpine as server-build

WORKDIR /build

# Install system dependencies
RUN apk update && apk add git gcc musl-dev

# Fetch and cache Go modules
COPY services/go.mod .
COPY services/go.sum .
RUN go mod download

# Copy in Go source files
COPY services/frontend-host/ ./service
COPY services/common/ ./common

# Now run the build
# Disabling cgo results in a fully static binary that can run without C libs
# Also inject version and build details 
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.version=$version -X 'main.buildInfo=$buildInfo'" \
    -o server ./service

# ================================================================================================
# === Stage 2: Build and bundle the Vue.js app with Vue CLI 3 ====================================
# ================================================================================================
FROM node:12-alpine as frontend-build
ARG sourceDir="web/client"

WORKDIR /build

# Install all the Vue.js dev tools & CLI, and our app dependencies 
COPY frontend/package*.json ./
RUN npm install --silent

# Copy in the Vue.js app source
#COPY frontend/.env.production .
COPY frontend/public ./public
COPY frontend/src ./src

# Run ESLint checks
RUN npm run lint
# Now main Vue CLI build & bundle, this will output to ./dist
RUN npm run build

# ================================================================================================
# === Stage 3: Bundle server exe and Vue dist in runtime image ===================================
# ================================================================================================
FROM scratch
WORKDIR /app 

# Copy in output from Vue bundle (the dist)
COPY --from=frontend-build /build/dist ./dist
# Copy the Go server binary
COPY --from=server-build /build/server . 

EXPOSE 8000

# That's it! Just run the server 
CMD [ "./server"]