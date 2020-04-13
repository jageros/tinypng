#plat ?= linux
#plats = linux darwin

all: cmg

define build_cmg
        @echo 'building $(1) ...'
        @GOOS=$(2) GOARCH=amd64 go build -o ./build/cmg ./$(1)
        @echo 'build $(1) done'
endef

cmg:
	$(call build_cmg,cmg,$(plat))

.PHONY: cmg