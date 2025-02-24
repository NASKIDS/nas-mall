v ?= $(shell git describe --tags --always --dirty)
ctx ?= docker-desktop

.PHONY: all
all: help

default: help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Initialize Project
.PHONY: init
init: ## Just copy `.env.example` to `.env` with one click, executed once.
	@scripts/copy_env.sh

##@ Build

.PHONY: gen-client
gen-client: ## gen client code of {svc}. example: make gen-client svc=product
	cd rpc_gen && cwgo client --type RPC --service ${svc} --module github.com/naskids/nas-mall/rpc_gen  -I ../idl  --idl ../idl/${svc}.proto

.PHONY: gen-server
gen-server: ## gen service code of {svc}. example: make gen-server svc=product
	cd app/${svc} && cwgo server --type RPC --service ${svc} --module github.com/naskids/nas-mall/app/${svc} --pass "-use github.com/naskids/nas-mall/rpc_gen/kitex_gen"  -I ../../idl  --idl ../../idl/${svc}.proto

.PHONY: gen
gen: gen-client gen-server ## gen client and svc code of {svc}. example: make gen svc=product

.PHONY: gen-frontend
gen-frontend: ## gen handler from one of proto files of frontend
	cd app/frontend && for proto_file in ../../idl/frontend/*.proto; \
    do \
        filename=$$(basename $$proto_file); \
        if [ $$filename = "common.proto" ]; then \
          continue; \
        fi; \
        cwgo server -I ../../idl --type HTTP --service frontend --module github.com/naskids/nas-mall/app/frontend --idl $$proto_file; \
    done

##@ Build

.PHONY: watch-frontend
watch-frontend:
	@cd app/frontend && air

.PHONY: tidy
tidy: ## run `go mod tidy` for all go module
	scripts/tidy.sh

.PHONY: lint
lint: ## run `gofmt` for all go module
	gofmt -l -w app
	gofumpt -l -w  app

.PHONY: vet
vet: ## run `go vet` for all go module
	scripts/vet.sh

.PHONY: lint-fix
lint-fix: ## run `golangci-lint` for all go module
	scripts/fix.sh

.PHONY: test
test: ## go unit test
# TODO go test

.PHONY: bin
bin: ## build binaries
	scripts/build_all.sh

.PHONY: run
run: ## run {svc} server. example: make run svc=product
	scripts/run.sh ${svc}

.PHONY: clean
clean: ## clean up all the tmp files
	rm -r app/**/log/ app/**/tmp/ app/**/output/ app/**/nohup.out

##@ Development Env With Docker Compose

.PHONY: env-start
env-start:  ## launch all middleware software as the docker
	@docker-compose up -d

.PHONY: env-stop
env-stop: ## stop all docker
	@docker-compose down

##@ Open Browser

.PHONY: open.gomall
open-gomall: ## open `gomall` website in the default browser
	@open "http://localhost:8080/"

.PHONY: open.consul
open-consul: ## open `consul ui` in the default browser
	@open "http://localhost:8500/ui/"

.PHONY: open.jaeger
open-jaeger: ## open `jaeger ui` in the default browser
	@open "http://localhost:16686/search"

.PHONY: open.prometheus
open-prometheus: ## open `prometheus ui` in the default browser
	@open "http://localhost:9090"

##@ Build Images

.PHONY: build-frontend
build-frontend: ## build the frontend service image
	docker build -f ./deploy/Dockerfile.frontend -t frontend:${v} .

.PHONY: build-svc
build-svc:  ## build one service image
	docker build -f ./deploy/Dockerfile.svc -t ${svc}:${v} --build-arg SVC=${svc} .

.PHONY: build-all
build-all: tidy vet lint-fix test ## build all service image
	docker build -f ./deploy/Dockerfile.frontend -t frontend:${v} -t frontend:latest .
	docker build -f ./deploy/Dockerfile.svc -t cart:${v} -t cart:latest --build-arg SVC=cart .
	docker build -f ./deploy/Dockerfile.svc -t checkout:${v} -t checkout:latest --build-arg SVC=checkout .
	docker build -f ./deploy/Dockerfile.svc -t email:${v} -t email:latest --build-arg SVC=email .
	docker build -f ./deploy/Dockerfile.svc -t order:${v} -t order:latest --build-arg SVC=order .
	docker build -f ./deploy/Dockerfile.svc -t payment:${v} -t payment:latest --build-arg SVC=payment .
	docker build -f ./deploy/Dockerfile.svc -t product:${v} -t product:latest --build-arg SVC=product .
	docker build -f ./deploy/Dockerfile.svc -t user:${v} -t user:latest --build-arg SVC=user .

##@ Deploy Images

.PHONY: deploy
deploy: ## deploy manifests to kubernetes
	kubectl apply --context=${ctx} -f deploy/gomall-dev-base.yaml
	kubectl apply --context=${ctx} -f deploy/gomall-dev-app.yaml