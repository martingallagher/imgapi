.PHONY: service

get_dep:
	command -v dep || go get -u github.com/golang/dep/cmd/dep

deps: get_dep
	dep ensure -vendor-only

update_deps: get_dep
	go get -u ./...
	rm -rf Gopkg.* vendor/
	dep init

service:
	protoc -I=service/ --go_out=plugins=grpc:service/ service/*.proto

check_go_fmt:
	@if [ -n "$$(gofmt -d $$(find . -name '*.go' -not -path './vendor/*'))" ]; then \
		>&2 echo "The .go sources aren't formatted. Please format them with 'go fmt'."; \
		exit 1; \
	fi

lint: check_go_fmt
	command -v gometalinter || go get -u github.com/alecthomas/gometalinter
	gometalinter --install >/dev/null
	gometalinter ./... --vendor --skip=vendor --exclude=\.*_mock\.*\.go --exclude=\.*pb\.*\.go --exclude=vendor\.* --cyclo-over=20 --deadline=2m --disable-all --enable-gc \
	--enable=errcheck \
	--enable=vet \
	--enable=deadcode \
	--enable=gocyclo \
	--enable=golint \
	--enable=varcheck \
	--enable=structcheck \
	--enable=maligned \
	--enable=ineffassign \
	--enable=interfacer \
	--enable=unconvert \
	--enable=goconst \
	--enable=gosimple \
	--enable=staticcheck \
	--enable=gas

test:
	go test -failfast -v -cover ./...

integration_test:
	go test -failfast -v -cover -tags=integration ./...

build:
	go build ./cmd/...
