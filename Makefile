export GOBIN=$(PWD)/gobin
OUT := bin/bus-routes

clean:
	rm -rf ./bin/*
.PHONY: clean

prepare.tools:
	cat tools.go | grep _ | grep -v "// _" | awk -F'"' '{print $$2}' | xargs -tI % go get %
.PHONY: prepare.tools

build: clean
	go build -o $(OUT) ./cmd
.PHONY: build

run: check build
	$(OUT) -config ./config.example.json
.PHONY: run

check: lint
.PHONY: check

lint:
	$(GOBIN)/golangci-lint run --timeout 3m
.PHONY: lint

check.swagger:
	$(GOBIN)/swagger validate assets/swagger/swagger.yml
.PHONY: check.swagger

gen.config:
	@if [ "$(ENV)" = "" ]; then echo "usage: make gen.config ENV=prod"; exit 1; fi
	@if [ "$(ENV)" = "prod" ]; then \
		echo "production is not ready"; exit 1; fi
	/c/windows/system32/wsl ansible-playbook -i provision/hosts -l $(ENV) provision/config-play.yml
.PHONY: gen.config

gen.config.local: ENV=local
gen.config.local: gen.config
.PHONY: gen.config.local

git.co:
	git checkout feature/FAT-15806/working-with-$(ph)
.PHONY: git.co

git.con:
	git checkout -b feature/FAT-15806/working-with-$(ph)
.PHONY: git.con

migration:
	go run ./scripts/migration_gen/migration_gen.go -name=$(name)