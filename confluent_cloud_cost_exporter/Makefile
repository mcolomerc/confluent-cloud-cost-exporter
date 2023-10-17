# Ensure that 'all' is the default target otherwise it will be the first target from Makefile.common.
all::

# Needs to be defined before including Makefile.common to auto-generate targets
DOCKER_ARCHS ?= amd64 armv7 arm64 ppc64le
DOCKER_REPO  ?= confluent-cloud-costs-exporter

include Makefile.common

DOCKER_IMAGE_NAME       ?= confluent-cloud-costs-exporter
