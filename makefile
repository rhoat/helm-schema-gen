reinstall:
	helm uninstall plugin schema-gen
	rm -rf $(helm env HELM_CACHE_HOME)
	helm plugin install https://github.com/rhoat/helm-schema-gen
install:
	helm plugin install https://github.com/rhoat/helm-schema-gen
install-local:
	helm plugin install .