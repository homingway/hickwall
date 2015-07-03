#
# Makefile
# oliveagle, 2015-07-02 14:12
#

.PHONY: all hickwall test clean

HICKWALL_BIN="hickwall"
HICKWALL_BIN_ARCH="hickwall"
HELPER_BIN="bin/hickwall_helper.exe"

# detect OS
ifeq ($(OS),Windows_NT)
    HICKWALL_BIN="hickwall.exe"
	HICKWALL_BIN_ARCH="hickwall-windows-386.exe"
	PACK=hickwall pack_win
endif

# add version and git commit hash
VER=$$(cat release-version)
GIT_HASH=$$(git rev-parse --short HEAD)
LD_FLAGS="-X main.Version $(VER) -X main.Build $(GIT_HASH)"

default: hickwall

all: $(PACK)

test:
	go test ./... -v | grep -E "(--- FAIL)|(^FAIL\s+)|(^ok\s+)"

hickwall: *.go
	rm -f $(HICKWALL_BIN)
	go build -ldflags $(LD_FLAGS) -v -o $(HICKWALL_BIN) && cp $(HICKWALL_BIN) bin/$(HICKWALL_BIN_ARCH)	

clean:
	rm -f $(HICKWALL_BIN)

# -------- windows ----------
win_helper:
	rm -f bin/hickwall_helper.exe
	gcc -Os make/win_helper_service/hickwall_helper.c -o bin/hickwall_helper.exe

pack_win: hickwall win_helper
	./make/pack_win/pack_with_inno.sh

# vim:ft=make
#

