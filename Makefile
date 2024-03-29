.DEFAULT_GOAL := help

.PHONY: help
help:	## https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

EXTERNAL_TOOLS := \
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1 \
	go install github.com/cosmtrek/air@latest \
	golang.org/x/pkgsite/cmd/pkgsite@latest # latest は go 1.19 以上が必要: https://github.com/golang/pkgsite#requirements

.PHONY: bootstrap
bootstrap: ## 外部ツールをインストールする。
	for t in $(EXTERNAL_TOOLS); do \
		echo "Installing $$t ..." ; \
		go install $$t ; \
	done

DC = docker compose

.PHONY: psql
psql:	## docker compose で起動した postgresql の db に接続する。
	$(DC) exec postgres psql -U root postgresql 

.PHONY: godoc
godoc:	## godoc をローカルで表示する。http://localhost:8080/{module_name}
	pkgsite

.PHONY: lint
lint:	## golangci を使って lint を走らせる。
	golangci-lint run -v

.PHONY: lint-fix
lint-fix:	## lint 実行時, gofumpt のエラーが出たらやると良い。
	golangci-lint run --fix

.PHONY: serve
serve:	## サーバーを起動する。
	go run -race app/*

.PHONY: dev
dev:	## Hot reload 付きでサーバーを起動する。
	air -c .air.toml

.PHONY: build-local
build-local:	## バイナリをビルドする（race オプションがついているため、ローカル実行専用とする）。
	go build -race -o app-local app/*

# カバレッジが低い場合は build-loacl でも動かしてみて競合の確認をしたい。
.PHONY: test
test:	## 全テストを実行する。
	go test -race -cover -shuffle=on ./... -v
