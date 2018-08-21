# Package configuration
PROJECT = lookout-gometalint-analyzer
COMMANDS = cmd/gometalint-analyzer

DOCKERFILES = Dockerfile:$(PROJECT)
DOCKER_ORG = "abezzubov"

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_BRANCH ?= v1
CI_PATH ?= .ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	git clone --quiet --depth 1 -b $(CI_BRANCH) $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

PROTOC := protoc
PROTOC_VER := "3.6.0"
# Generate go code from proto files
.PHONY: check-protoc
check-protoc:
	./_tools/install-protoc-maybe.sh
.PHONY: protogen
protogen: check-protoc
	$(GOCMD) install ./vendor/github.com/gogo/protobuf/protoc-gen-gogofaster
	$(PROTOC) \
		-I sdk \
		--gogofaster_out=plugins=grpc,\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:pb \
sdk/*.proto
