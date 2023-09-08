PROJECT_NAME := Pulumi Command Resource Provider

PACK             := command
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-command
NODE_MODULE_NAME := @pulumi/command
NUGET_PKG_NAME   := Pulumi.Command

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         ?= $(shell pulumictl get version)
PROVIDER_PATH   := provider
VERSION_PATH    := ${PROVIDER_PATH}/pkg/version.Version

SCHEMA_FILE     := provider/cmd/pulumi-resource-command/schema.json
GOPATH          := $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

# Need to pick up locally pinned pulumi-langage-* plugins.
export PATH := .pulumi/bin:$(PATH)

ensure:: tidy

tidy:
	cd provider && go mod tidy && cd tests && go mod tidy
	cd sdk && go mod tidy
	cd examples && go mod tidy

codegen::
	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && go build -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))
	$(WORKING_DIR)/bin/${CODEGEN} $(SCHEMA_FILE) --version ${VERSION} 

provider::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider_debug::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd provider/pkg && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...

dotnet_sdk:: DOTNET_VERSION := $(shell pulumictl get version --language dotnet)
dotnet_sdk::	.pulumi/bin/pulumi
	rm -rf sdk/dotnet
	.pulumi/bin/pulumi package gen-sdk --language dotnet $(SCHEMA_FILE)
	# Copy the logo to the dotnet directory before building so it can be included in the nuget package archive.
	# https://github.com/pulumi/pulumi-command/issues/243
	cd ${PACKDIR}/dotnet/&& \
		echo "${DOTNET_VERSION}" >version.txt && \
		cp $(WORKING_DIR)/assets/logo.png logo.png && \
		dotnet build /p:Version=${DOTNET_VERSION}

go_sdk::	.pulumi/bin/pulumi
	rm -rf sdk/go
	.pulumi/bin/pulumi package gen-sdk --language go $(SCHEMA_FILE)

nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk::	.pulumi/bin/pulumi
	rm -rf sdk/nodejs
	.pulumi/bin/pulumi package gen-sdk --language nodejs $(SCHEMA_FILE)
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

python_sdk:: PYPI_VERSION := $(shell pulumictl get version --language python)
python_sdk::	.pulumi/bin/pulumi
	rm -rf sdk/python
	.pulumi/bin/pulumi package gen-sdk --language python $(SCHEMA_FILE)
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(PYPI_VERSION)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist

bin/pulumi-java-gen::
	echo pulumi-java-gen is no longer necessary

java_sdk:: PACKAGE_VERSION := $(shell pulumictl get version --language generic)
java_sdk::	.pulumi/bin/pulumi
	rm -rf sdk/java
	.pulumi/bin/pulumi package gen-sdk --language java $(SCHEMA_FILE)
	cd sdk/java/ && \
		gradle --console=plain build

.PHONY: build
build:: provider build_sdks

.PHONY: build_sdks
build_sdks: codegen dotnet_sdk go_sdk nodejs_sdk python_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run --timeout 10m && popd ; \
	done


install:: install_nodejs_sdk install_dotnet_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin


GO_TEST := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

test_all::
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
	cd tests/sdk/go && $(GO_TEST) ./...

install_dotnet_sdk::
	rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	mkdir -p $(WORKING_DIR)/nuget
	find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

install_python_sdk::
	#target intentionally blank

install_go_sdk::
	#target intentionally blank

install_java_sdk::
	#target intentionally blank

install_nodejs_sdk::
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

test_unit: tidy
	cd provider/tests && go test -v ./...

test:: tidy test_unit
	cd examples && go test -v -tags=all -timeout 2h

# --------- File-based targets --------- #

.pulumi/bin/pulumi: PULUMI_VERSION := $(shell cat .pulumi.version)
.pulumi/bin/pulumi: HOME := $(WORKING_DIR)
.pulumi/bin/pulumi: .pulumi.version
	curl -fsSL https://get.pulumi.com | sh -s -- --version "$(PULUMI_VERSION)"
