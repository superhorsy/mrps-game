.PHONY: goimports
goimports:
	goimports -local "mrps-game" -w internal
	goimports -local "mrps-game" -w cmd

.PHONY: lint-local
lint-local: ## Run lint locally
	golangci-lint run -v --deadline 300s