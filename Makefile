PKG = github.com/k1LoW/frgm
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: depsdev test sec

test:
	go test ./... -coverprofile=coverage.out -covermode=count

sec:
	gosec ./...

build:
	packr2
	go build -ldflags="$(BUILD_LDFLAGS)"
	packr2 clean

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@v0.10.2
	go install github.com/Songmu/gocredits/cmd/gocredits@v0.2.0
	go install github.com/securego/gosec/v2/cmd/gosec@v2.8.1
	go install github.com/gobuffalo/packr/v2/packr2@v2.8.3

prerelease:
	git push origin main --tag
	ghch -w -N ${VER}
	gocredits -skip-missing . > CREDITS
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS
	git commit -m'Bump up version number'
	git tag ${VER}

release:
	git push origin main --tag
	goreleaser --rm-dist

.PHONY: default test
