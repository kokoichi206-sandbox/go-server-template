.PHONY: build help bootstrap godoc
.DEFAULT_GOAL := help

EXTERNAL_TOOLS := \
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1 \
	golang.org/x/pkgsite/cmd/pkgsite@latest # latest は go 1.19 以上が必要: https://github.com/golang/pkgsite#requirements

help:	## https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

bootstrap: ## 外部ツールをインストールする。
	for t in $(EXTERNAL_TOOLS); do \
		echo "Installing $$t ..." ; \
		go install $$t ; \
	done

DC = docker compose
psql:	## docker compose で起動した postgresql の db に接続する。
	$(DC) exec postgresql psql -U root postgresql 

godoc:	## godoc をローカルで表示する。http://localhost:8080/{module_name}
	pkgsite

.PHONY: lint lint-fix serve

lint:	## golangci を使って lint を走らせる。
	golangci-lint run -v

lint-fix:	## lint 実行時, gofumpt のエラーが出たらやると良い。
	golangci-lint run --fix

serve:	## サーバーを起動する。
	go run app/*

test:	## 全テストを実行する。
	go test -cover -shuffle=on ./... -v
