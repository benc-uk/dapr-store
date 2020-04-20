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
VERSION ?= 0.0.1

# Most likely want to override these when calling `make docker`
DOCKER_REG ?= docker.io
DOCKER_REPO ?= daprstore
DOCKER_TAG ?= latest
DOCKER_PREFIX := $(DOCKER_REG)/$(DOCKER_REPO)

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
	echo "BLEEP BLOOP. All tests passing. Nothing to see here. Move along"

################################################################################
# Gofmt
################################################################################
.PHONY: gofmt
gofmt :
	@./.github/workflows/gofmt-action.sh $(SERVICE_DIR)/

################################################################################
# Build Docker images
################################################################################
.PHONY: docker
docker :
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg SERVICE_NAME=cart \
	--build-arg SERVICE_PORT=9001 \
	-t $(DOCKER_PREFIX)/cart:$(DOCKER_TAG)
	
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg SERVICE_NAME=products \
	--build-arg SERVICE_PORT=9002 \
	--build-arg CGO_ENABLED=1 \
	-t $(DOCKER_PREFIX)/products:$(DOCKER_TAG)

	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg SERVICE_NAME=users \
	--build-arg SERVICE_PORT=9003 \
	-t $(DOCKER_PREFIX)/users:$(DOCKER_TAG)

	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg SERVICE_NAME=orders \
	--build-arg SERVICE_PORT=9004 \
	-t $(DOCKER_PREFIX)/orders:$(DOCKER_TAG)

	docker build . -f build/frontend.Dockerfile \
	--build-arg VERSION=$(VERSION) \
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
