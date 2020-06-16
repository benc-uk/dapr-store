# ------------------------------------------------------------
# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.
# ------------------------------------------------------------

################################################################################
# Variables
################################################################################
CGO := 0
SERVICE_DIR := cmd
FRONTEND_DIR := web/frontend
OUTPUT_DIR := ./output
VERSION ?= 0.0.1
BUILD_INFO ?= "Makefile build"

# Most likely want to override these when calling `make docker`
DOCKER_REG ?= docker.io
DOCKER_REPO ?= daprstore
DOCKER_TAG ?= latest
DOCKER_PREFIX := $(DOCKER_REG)/$(DOCKER_REPO)

# Change if you know what you're doing
CLIENT_ID ?= ""

################################################################################
# Lint check everything
################################################################################
.PHONY: lint
lint : $(FRONTEND_DIR)/node_modules
	golint -set_exit_status $(SERVICE_DIR)/...
	@cd $(FRONTEND_DIR); npm run lint


################################################################################
# Run tests
################################################################################
.PHONY: test
test : 
	go test -v github.com/benc-uk/dapr-store/cmd/cart
	go test -v github.com/benc-uk/dapr-store/cmd/orders
	go test -v github.com/benc-uk/dapr-store/cmd/products
	go test -v github.com/benc-uk/dapr-store/cmd/users
	@cd $(FRONTEND_DIR); NODE_ENV=test npm run test -- --ci


################################################################################
# Run tests with output
################################################################################
.PHONY: test-output
test-output : 
	rm -rf $(OUTPUT_DIR) && mkdir -p $(OUTPUT_DIR)
	gotestsum --junitfile $(OUTPUT_DIR)/unit-tests.xml ./cmd/cart ./cmd/users ./cmd/products ./cmd/orders --coverprofile $(OUTPUT_DIR)/coverage
	@cd $(FRONTEND_DIR); NODE_ENV=test npm run test -- --ci


################################################################################
# Prepare HTML reports from test output
################################################################################
.PHONY: reports
reports : 
	./web/frontend/node_modules/xunit-viewer/bin/xunit-viewer -r $(OUTPUT_DIR)/unit-tests.xml -o $(OUTPUT_DIR)/unit-tests.html
	./web/frontend/node_modules/xunit-viewer/bin/xunit-viewer -r $(OUTPUT_DIR)/unit-tests-frontend.xml -o $(OUTPUT_DIR)/unit-tests-frontend.html
	go tool cover -html=$(OUTPUT_DIR)/coverage -o $(OUTPUT_DIR)/cover.html
	cp testing/reports.html $(OUTPUT_DIR)/index.html


################################################################################
# Gofmt
################################################################################
.PHONY: gofmt
gofmt :
	@./.github/workflows/gofmt-action.sh $(SERVICE_DIR)/


################################################################################
# Clean up project
################################################################################
.PHONY: clean
clean :
	rm -rf $(FRONTEND_DIR)/node_modules
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/coverage
	rm -rf output
	rm -rf $(SERVICE_DIR)/cart/cart
	rm -rf $(SERVICE_DIR)/orders/orders
	rm -rf $(SERVICE_DIR)/users/users
	rm -rf $(SERVICE_DIR)/products/products
	rm -rf $(SERVICE_DIR)/frontend-host/frontend-host


################################################################################
# Build Docker images
################################################################################
.PHONY: docker
docker :
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO=$(BUILD_INFO) \
	--build-arg SERVICE_NAME=cart \
	--build-arg SERVICE_PORT=9001 \
	-t $(DOCKER_PREFIX)/cart:$(DOCKER_TAG)
	
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO=$(BUILD_INFO) \
	--build-arg SERVICE_NAME=products \
	--build-arg SERVICE_PORT=9002 \
	--build-arg CGO_ENABLED=1 \
	-t $(DOCKER_PREFIX)/products:$(DOCKER_TAG)

	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO=$(BUILD_INFO) \
	--build-arg SERVICE_NAME=users \
	--build-arg SERVICE_PORT=9003 \
	-t $(DOCKER_PREFIX)/users:$(DOCKER_TAG)

	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO=$(BUILD_INFO) \
	--build-arg SERVICE_NAME=orders \
	--build-arg SERVICE_PORT=9004 \
	-t $(DOCKER_PREFIX)/orders:$(DOCKER_TAG)

	docker build . -f build/frontend.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg CLIENT_ID=$(CLIENT_ID) \
	-t $(DOCKER_PREFIX)/frontend-host:$(DOCKER_TAG)


################################################################################
# Build Docker image for frontend only
################################################################################
.PHONY: docker-frontend
docker-frontend :
	docker build . -f build/frontend.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg CLIENT_ID=$(CLIENT_ID) \
	-t $(DOCKER_PREFIX)/frontend-host:$(DOCKER_TAG)


################################################################################
# Push Docker images
################################################################################
.PHONY: push
push :
	docker push $(DOCKER_PREFIX)/cart:$(DOCKER_TAG)
	docker push $(DOCKER_PREFIX)/products:$(DOCKER_TAG)
	docker push $(DOCKER_PREFIX)/users:$(DOCKER_TAG)
	docker push $(DOCKER_PREFIX)/orders:$(DOCKER_TAG)
	docker push $(DOCKER_PREFIX)/frontend-host:$(DOCKER_TAG)
	

################################################################################
# Frontend / Vue.js
################################################################################
frontend : $(FRONTEND_DIR)/node_modules
	cd $(FRONTEND_DIR); npm run build
	cd $(SERVICE_DIR)/frontend-host; go build

$(FRONTEND_DIR)/node_modules: $(FRONTEND_DIR)/package.json
	cd $(FRONTEND_DIR); npm install --silent
	touch -m $(FRONTEND_DIR)/node_modules

$(FRONTEND_DIR)/package.json: 
	@echo "package.json was modified"
