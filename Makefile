tag: VERSION = $(shell grep "const CLIVersion" main.go | awk 'NF>1{print $$NF}' | cut -d '"' -f2)
tag:
	git tag -a $(VERSION)