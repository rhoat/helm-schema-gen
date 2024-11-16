VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
PKG:= github.com/rhoat/helm-schema-gen
LDFLAGS := -X $(PKG)/pkg/commands/version.Version=$(VERSION)
HELM_3_PLUGINS := $(shell helm env HELM_PLUGINS)
reinstall:
	helm uninstall plugin schema-gen
	rm -rf $(helm env HELM_CACHE_HOME)
	helm plugin install https://github.com/rhoat/helm-schema-gen
install:
	helm plugin install https://github.com/rhoat/helm-schema-gen
install-local:
	helm plugin install .

.PHONY: build
build:
	mkdir -p bin/
	go build -v -o bin/schema-gen -ldflags="$(LDFLAGS)"

.PHONY: releaser
releaser:
	CGO_ENABLED=0 VERSION=$(VERSION) goreleaser release --skip=publish --clean

.PHONY: install/helm3
install/helm3: build
	mkdir -p $(HELM_3_PLUGINS)/helm-schema-gen/bin
	cp bin/schema-gen $(HELM_3_PLUGINS)/helm-schema-gen/bin
	cp plugin.yaml $(HELM_3_PLUGINS)/helm-schema-gen/

.PHONY: tag
	git push origin :refs/tags/$(VERSION)
	git tag $(VERSION)
	git push origin --tags