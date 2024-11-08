GO ?= go

# Override in Makefile the working directory.
PWD ?= $(shell pwd)

MAKEFILES_PATH ?= $(PWD)/pkg/makefiles

-include $(MAKEFILES_PATH)/help.mk
