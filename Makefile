SERVICE_DIR := cmd
FRONTEND_DIR := web/frontend
OUTPUT_DIR := ./output
VERSION ?= 0.8.3
BUILD_INFO ?= "Local makefile build"
DAPR_RUN_LOGLEVEL := warn

# Most likely want to override these when calling `make image-all`
IMAGE_REG ?= ghcr.io
IMAGE_REPO ?= benc-uk/daprstore
IMAGE_TAG ?= latest
IMAGE_PREFIX := $(IMAGE_REG)/$(IMAGE_REPO)

.EXPORT_ALL_VARIABLES:
.PHONY: help lint lint-fix test test-reports docker-build docker-run docker-stop docker-push bundle clean run stop
.DEFAULT_GOAL := help

help:  ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: $(FRONTEND_DIR)/node_modules      ## 🔎 Lint & format, check to be run in CI, sets exit code on error
	golangci-lint run --modules-download-mode=mod --timeout=4m ./...
	cd $(FRONTEND_DIR); npm run lint

lint-fix: $(FRONTEND_DIR)/node_modules  ## 📝 Lint & format, fixes errors and modifies code
	golangci-lint run --modules-download-mode=mod --timeout=4m --fix ./...
	cd $(FRONTEND_DIR); npm run lint-fix

test:  ## 🎯 Unit tests for services and snapshot tests for SPA frontend 
	go test -v -count=1 ./$(SERVICE_DIR)/...
	@cd $(FRONTEND_DIR); NODE_ENV=test npm run test -- --ci

test-reports: $(FRONTEND_DIR)/node_modules  ## 📜 Unit tests with coverage and test reports (deprecated)
	@rm -rf $(OUTPUT_DIR) && mkdir -p $(OUTPUT_DIR)
	@which gotestsum || go get gotest.tools/gotestsum
	gotestsum --junitfile $(OUTPUT_DIR)/unit-tests.xml ./$(SERVICE_DIR)/... --coverprofile $(OUTPUT_DIR)/coverage
	cd $(FRONTEND_DIR); NODE_ENV=test npm run test -- --ci
	./$(FRONTEND_DIR)/node_modules/xunit-viewer/bin/xunit-viewer -r $(OUTPUT_DIR)/unit-tests.xml -o $(OUTPUT_DIR)/unit-tests.html
	./$(FRONTEND_DIR)/node_modules/xunit-viewer/bin/xunit-viewer -r $(OUTPUT_DIR)/unit-tests-frontend.xml -o $(OUTPUT_DIR)/unit-tests-frontend.html
	go tool cover -html=$(OUTPUT_DIR)/coverage -o $(OUTPUT_DIR)/cover.html
	cp testing/reports.html $(OUTPUT_DIR)/index.html
	
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

run:  ## 🚀 Start & run everything locally as processes
	cd $(FRONTEND_DIR); npm run serve &
	dapr run --app-id cart     --app-port 9001 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/cart &
	dapr run --app-id products --app-port 9002 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/products ./cmd/products/sqlite.db &
	dapr run --app-id users    --app-port 9003 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/users &
	dapr run --app-id orders   --app-port 9004 --log-level $(DAPR_RUN_LOGLEVEL) go run github.com/benc-uk/dapr-store/cmd/orders &
	@sleep 6
	@./scripts/local-gateway/run.sh &
	@sleep infinity
	@echo "!!! Processes may still be running, please run `make stop` in order to shutdown everything"

docker-run: ## 🐋 Run locally using containers and Docker compose
	@./scripts/local-gateway/run.sh &
	@docker compose -f ./build/compose.yaml up --remove-orphans

docker-build: ## 🔨 Build all containers using Docker compose
	docker compose -f ./build/compose.yaml build

docker-push: ## 📤 Push all containers using Docker compose
	docker compose -f ./build/compose.yaml push

docker-stop: ## 🚫 Stop and remove local containers
	docker rm -f api-gateway || true
	docker compose -f ./build/compose.yaml rm -f

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

# ===============================================================================

docker-build-cart:
	docker compose -f ./build/compose.yaml build cart

docker-build-products:
	docker compose -f ./build/compose.yaml build products

docker-build-users:
	docker compose -f ./build/compose.yaml build users

docker-build-orders:
	docker compose -f ./build/compose.yaml build orders

docker-build-frontend:
	docker compose -f ./build/compose.yaml build frontend

# ===============================================================================

docker-push-cart:
	docker compose -f ./build/compose.yaml push cart

docker-push-products:
	docker compose -f ./build/compose.yaml push products

docker-push-users:
	docker compose -f ./build/compose.yaml push users

docker-push-orders:
	docker compose -f ./build/compose.yaml push orders

docker-push-frontend:
	docker compose -f ./build/compose.yaml push frontend

