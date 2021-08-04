plat ?= darwin
plats = linux darwin

all: cmg

define build_cmg
        @echo 'building $(1) ...'
        @GOOS=$(2) GOARCH=amd64 go build -o ./builder/cmg ./$(1)
        @echo 'build $(1) done'
endef

cmg:
	$(call build_cmg,cmg,$(plat))

.PHONY: cmg