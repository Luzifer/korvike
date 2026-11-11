default:

korvike:
	go build \
		-buildvcs=false \
		-ldflags "-s -w -buildid=" \
		-mod=readonly \
		-trimpath \
		-o korvike

publish:
	bash ./ci/build.sh

test:
	cd functions && go test -cover -v
	golangci-lint run ./...

README.md: export KORVIKE=$(CURDIR)/korvike
README.md: korvike
	find . -name generate.go -execdir env GOWORK=off go run . \; >$@

trivy:
	trivy fs . \
		--dependency-tree \
		--exit-code 1 \
		--format table \
		--ignore-unfixed \
		--quiet \
		--scanners misconfig,license,secret,vuln \
		--severity HIGH,CRITICAL

.PHONY: README.md korvike
