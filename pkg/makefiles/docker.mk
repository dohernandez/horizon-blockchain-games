
# Override in app Makefile to control docker file path.
DOCKERFILE_PATH ?= Dockerfile

# Override in app Makefile to control docker build context.
DOCKERBUILD_CONTEXT ?= .

# Override in app Makefile to control docker image tag.
DOCKER_IMAGE_TAG ?= latest

# Override in app Makefile to control docker image github token in case the docker required it into build.
DOCKER_GITHUB_TOKEN ?= ""

# Override in app Makefile to control docker docker-compose.yml build using secret instead of args.
DOCKER_SECRET ?= false

## Build docker image
build-image:
	@DOCKER_IMAGE_TAG=$(DOCKER_IMAGE_TAG) \
	DOCKERFILE_PATH=$(DOCKERFILE_PATH) \
	DOCKERBUILD_CONTEXT=$(DOCKERBUILD_CONTEXT) \
	DOCKER_GITHUB_TOKEN=$(DOCKER_GITHUB_TOKEN) \
	DOCKER_SECRET=$(DOCKER_SECRET) \
	bash $(MAKEFILES_PATH)/docker-build.sh


.PHONY: build-image