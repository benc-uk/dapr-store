SERVICE_DIR := cmd
FRONTEND_DIR := web/frontend
OUTPUT_DIR := ./output
VERSION ?= 0.6.0
BUILD_INFO ?= "Makefile build"
DAPR_RUN_LOGLEVEL := warn

REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint

# Most likely want to override these when calling `make image-all`
IMAGE_REG ?= ghcr.io
IMAGE_REPO ?= benc-uk/daprstore
IMAGE_TAG ?= latest
IMAGE_PREFIX := $(IMAGE_REG)/$(IMAGE_REPO)
IMAGE_LIST := cart orders users products frontend

.PHONY: help lint lint-fix test test-reports test-snapshot image-all bundle clean run
.DEFAULT_GOAL := help

help:  ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: $(FRONTEND_DIR)/node_modules      ## 🔎 Lint & format, check to be run in CI, sets exit code on error 
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	cd $(SERVICE_DIR); $(GOLINT_PATH) run --modules-download-mode=mod ./...
	cd $(FRONTEND_DIR); npm run lint

lint-fix: $(FRONTEND_DIR)/node_modules  ## 📝 Lint & format, fixes errors and modifies code
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	cd $(SERVICE_DIR); $(GOLINT_PATH) run --modules-download-mode=mod ./... --fix
	cd $(FRONTEND_DIR); npm run lint-fix

test:  ## 🎯 Unit tests for services and snapshot tests for SPA frontend 
	go test -v ./$(SERVICE_DIR)/...
	@cd $(FRONTEND_DIR); NODE_ENV=test npm run test -- --ci

test-reports: $(FRONTEND_DIR)/node_modules  ## 📜 Unit tests with coverage and test reports
	@rm -rf $(OUTPUT_DIR) && mkdir -p $(OUTPUT_DIR)
	@which gotestsum || go get gotest.tools/gotestsum
	gotestsum --junitfile $(OUTPUT_DIR)/unit-tests.xml ./$(SERVICE_DIR)/... --coverprofile $(OUTPUT_DIR)/coverage
	cd $(FRONTEND_DIR); NODE_ENV=test npm run test -- --ci
	./$(FRONTEND_DIR)/node_modules/xunit-viewer/bin/xunit-viewer -r $(OUTPUT_DIR)/unit-tests.xml -o $(OUTPUT_DIR)/unit-tests.html
	./$(FRONTEND_DIR)/node_modules/xunit-viewer/bin/xunit-viewer -r $(OUTPUT_DIR)/unit-tests-frontend.xml -o $(OUTPUT_DIR)/unit-tests-frontend.html
	go tool cover -html=$(OUTPUT_DIR)/coverage -o $(OUTPUT_DIR)/cover.html
	cp testing/reports.html $(OUTPUT_DIR)/index.html

test-snapshot:  ## 📷 Update snapshots for frontend tests
	@cd $(FRONTEND_DIR); NODE_ENV=test npm run test-update

image-all:      ## 📦 Build all container images
	for img in $(IMAGE_LIST); do \
		make image-$$img ; \
	done

image-cart:
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO='$(BUILD_INFO)' \
	--build-arg SERVICE_NAME=cart \
	--build-arg SERVICE_PORT=9001 \
	-t $(IMAGE_PREFIX)/cart:$(IMAGE_TAG)

image-products:
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO='$(BUILD_INFO)' \
	--build-arg SERVICE_NAME=products \
	--build-arg SERVICE_PORT=9002 \
	--build-arg CGO_ENABLED=1 \
	-t $(IMAGE_PREFIX)/products:$(IMAGE_TAG)

image-users:
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO='$(BUILD_INFO)' \
	--build-arg SERVICE_NAME=users \
	--build-arg SERVICE_PORT=9003 \
	-t $(IMAGE_PREFIX)/users:$(IMAGE_TAG)

image-orders:
	docker build . -f build/service.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO='$(BUILD_INFO)' \
	--build-arg SERVICE_NAME=orders \
	--build-arg SERVICE_PORT=9004 \
	-t $(IMAGE_PREFIX)/orders:$(IMAGE_TAG)

image-frontend:
	docker build . -f build/frontend.Dockerfile \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_INFO='$(BUILD_INFO)' \
	-t $(IMAGE_PREFIX)/frontend-host:$(IMAGE_TAG)

push-all:  ## 📤 Push all images to registry 
	docker push $(IMAGE_PREFIX)/cart:$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/products:$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/users:$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/orders:$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/frontend-host:$(IMAGE_TAG)
	
bundle: $(FRONTEND_DIR)/node_modules  ## 💻 Build and bundle the frontend Vue SPA
	cd $(FRONTEND_DIR); npm run build
	cd $(SERVICE_DIR)/frontend-host; go build

clean:  ## 🧹 Clean the project, remove modules, binaries and outputs
	rm -rf output
	rm -rf $(FRONTEND_DIR)/node_modules
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/coverage
	rm -rf $(SERVICE_DIR)/cart/cart
	rm -rf $(SERVICE_DIR)/orders/orders
	rm -rf $(SERVICE_DIR)/users/users
	rm -rf $(SERVICE_DIR)/products/products
	rm -rf $(SERVICE_DIR)/frontend-host/frontend-host

run:  ## 🚀 Start & run everything locally
	cd $(FRONTEND_DIR); npm run serve &
	dapr run --app-id cart     --app-port 9001 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/cart &
	dapr run --app-id products --app-port 9002 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/products ./cmd/products/sqlite.db &
	dapr run --app-id users    --app-port 9003 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/users &
	dapr run --app-id orders   --app-port 9004 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/orders &
	@sleep 6
	@./scripts/local-gateway/run.sh &
	@sleep infinity
	@echo "!!! Processes may still be running, please run `make stop` in order to shutdown everything"

stop: ## ⛔ Stop & kill everything started locally from `make run`
	docker rm -f api-gateway || true
	dapr stop --app-id api-gateway
	dapr stop --app-id cart
	dapr stop --app-id products
	dapr stop --app-id users
	dapr stop --app-id orders
	pkill cart; pkill users; pkill orders; pkill products; pkill main

# ===============================================================================

$(FRONTEND_DIR)/node_modules: $(FRONTEND_DIR)/package.json
	cd $(FRONTEND_DIR); npm install --silent
	touch -m $(FRONTEND_DIR)/node_modules

$(FRONTEND_DIR)/package.json: 
	@echo "package.json was modified"
