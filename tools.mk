TOOLS_BIN_DIR := hack/tools/bin
export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

# Tool targets should declare go.mod as a prerequisite, if the tool's version is managed via go modules. This causes
# make to rebuild the tool in the desired version, when go.mod is changed.
# For tools where the version is not managed via go.mod, we use a file per tool and version as an indicator for make
# whether we need to install the tool or a different version of the tool (make doesn't rerun the rule if the rule is
# changed).

# Use this "function" to add the version file as a prerequisite for the tool target: e.g.
#   $(KUBECTL): $(call tool_version_file,$(KUBECTL),$(KUBECTL_VERSION))
tool_version_file = $(TOOLS_BIN_DIR)/.version_$(subst $(TOOLS_BIN_DIR)/,,$(1))_$(2)

# This target cleans up any previous version files for the given tool and creates the given version file.
# This way, we can generically determine, which version was installed without calling each and every binary explicitly.
$(TOOLS_BIN_DIR)/.version_%:
	@version_file=$@; rm -f $${version_file%_*}*
	@touch $@

GOIMPORTS_REVISER := $(TOOLS_BIN_DIR)/goimports-reviser
# renovate: datasource=github-releases depName=incu6us/goimports-reviser
GOIMPORTS_REVISER_VERSION ?= v3.3.1
$(GOIMPORTS_REVISER): $(call tool_version_file,$(GOIMPORTS_REVISER),$(GOIMPORTS_REVISER_VERSION))
	GOBIN=$(abspath $(TOOLS_BIN_DIR)) go install github.com/incu6us/goimports-reviser/v3@$(GOIMPORTS_REVISER_VERSION)

GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.54.2
$(GOLANGCI_LINT): $(call tool_version_file,$(GOLANGCI_LINT),$(GOLANGCI_LINT_VERSION))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_BIN_DIR) $(GOLANGCI_LINT_VERSION)

GORELEASER := $(TOOLS_BIN_DIR)/goreleaser
# renovate: datasource=github-releases depName=goreleaser/goreleaser
GORELEASER_VERSION ?= v1.26.1
$(GORELEASER): $(call tool_version_file,$(GORELEASER),$(GORELEASER_VERSION))
	curl -sSfL https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASER_VERSION)/goreleaser_$(shell uname -s)_$(shell uname -m).tar.gz | tar -xzmf - -C $(TOOLS_BIN_DIR) goreleaser
