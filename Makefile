CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
COMMIT=`git rev-parse --short HEAD`
APP?=quorra
REPO?=factionlabs/$(APP)
TAG?=latest
MEDIA_SRCS=$(shell find public/ -type f \
	-not -path "public/dist/*" \
	-not -path "public/node_modules/*")

export GOPATH:=$(PWD)/vendor:$(GOPATH)

all: media build image

add-deps:
	@godep save
	@rm -rf Godeps

build:
	@cd cmd/$(APP) && go build -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT)" .

build-static:
	@cd cmd/$(APP) && go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT)" .

dev-setup:
	@echo "This could take a while..."
	@npm install -g gulp browserify babelify
	@cd public && npm install

media: media-semantic media-app

media-semantic: public/dist/.bundle_timestamp
public/dist/.bundle_timestamp: $(MEDIA_SRCS)
	@cp -f public/semantic.theme.config public/semantic/src/theme.config
	@cp -r public/semantic.theme public/semantic/src/themes/app
	@cd public/semantic && gulp build
	@mkdir -p public/dist
	@cd public && rm -rf dist/semantic* dist/themes
	@cp -f public/semantic/dist/semantic.min.css public/dist/semantic.min.css
	@cp -f public/semantic/dist/semantic.min.js public/dist/semantic.min.js
	@mkdir -p public/dist/themes/default && cp -r public/semantic/dist/themes/default/assets public/dist/themes/default/
	@touch public/dist/.bundle_timestamp

media-app:
	@mkdir -p public/dist
	@cd public && rm -rf dist/bundle.js
	@# add frontend ui components here
	@cd public/src && browserify app/* dashboard/* -t babelify --outfile ../dist/bundle.js

image: build-static
	@mkdir -p build
	@cp -r cmd/$(APP)/$(APP) build/
	@cp -r public build/
	@rm -rf build/public/node_modules build/public/semantic/{gulpfile.js,src,tasks} build/public/semantic.json
	@docker build -t $(REPO):$(TAG) .

release: image
	@docker push $(REPO):$(TAG)

package:
	@mkdir -p build
	@cp -r cmd/$(APP)/$(APP) build/
	@cp -r public build/
	@cd build/public && rm -rf node_modules package.json semantic semantic.theme semantic.theme.config src

test: build
	@go test -v ./...

clean:
	@rm cmd/$(APP)/$(APP)
	@rm .build_timestamp
	@rm -rf build
	@rm -rf public/dist/*

.PHONY: add-deps build build-static dev-setup media media-semantic media-app image release test clean
