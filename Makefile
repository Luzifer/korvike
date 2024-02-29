publish:
	bash ./ci/build.sh

test:
	cd functions && go test -cover -v
	golangci-lint run ./...

trivy:
	trivy fs . \
		--dependency-tree \
		--exit-code 1 \
		--format table \
		--ignore-unfixed \
		--quiet \
		--scanners misconfig,license,secret,vuln \
		--severity HIGH,CRITICAL
