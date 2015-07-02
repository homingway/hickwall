#
# Makefile
# oliveagle, 2015-07-02 14:12
#


HICKWALL_BIN="hickwall"
HICKWALL_BIN_ARCH="hickwall"

# detect OS
ifeq ($(OS),Windows_NT)
    HICKWALL_BIN="hickwall.exe"
	HICKWALL_BIN_ARCH="hickwall-windows-386.exe"
endif

VER=$$(cat release-version)
GIT_HASH=$$(git rev-parse --short HEAD)
LD_FLAGS="-X main.Version $(VER) -X main.Build $(GIT_HASH)"

all:
	rm -f $(HICKWALL_BIN)
	go build -ldflags $(LD_FLAGS) -v -o $(HICKWALL_BIN) && cp $(HICKWALL_BIN) bin/$(HICKWALL_BIN_ARCH)


test:
	go test ./... -v | grep -E "(--- FAIL)|(^FAIL\s+)|(^ok\s+)"

# vim:ft=make
#

