GO ?= go
CMAKE ?= cmake
GO-SRC-DIR = src
GO-BUILD-FLAGS = -trimpath -mod=vendor
GO-LDFLAGS = -s -w
GO-BUILD-ENV = GO111MODULE=on GOPROXY=off GOTOOLCHAIN=local
GO-PKG-PREFIX = .
BUILD-DEST = build/release
GUI-BUILD-DIR = build/gui

VERSION ?= $(shell sed -n '1s/^[^(]*(\([^)]*\)).*/\1/p' debian/changelog)

.PHONY: build check clean gui gui-check inst-gui vendor

all: build install

clean:
	rm -rf build

$(BUILD-DEST):
	mkdir -p $(BUILD-DEST)

# Regenerate src/vendor after changing src/go.mod or src/go.sum.
# CI verifies vendor is in sync: run this and commit the result.
vendor:
	cd $(GO-SRC-DIR) && $(GO) mod vendor && $(GO) mod verify

anonimax-cli: $(BUILD-DEST)
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) $(GO) build $(GO-BUILD-FLAGS) -ldflags="$(GO-LDFLAGS) -X main.version=$(VERSION)" -o ../$(BUILD-DEST)/anonimax $(GO-PKG-PREFIX)/cmd/anonsurf-cli

dnstool: $(BUILD-DEST)
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) $(GO) build $(GO-BUILD-FLAGS) -ldflags="$(GO-LDFLAGS)" -o ../$(BUILD-DEST)/dnstool $(GO-PKG-PREFIX)/cmd/dnstool

make-torrc: $(BUILD-DEST)
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) $(GO) build $(GO-BUILD-FLAGS) -ldflags="$(GO-LDFLAGS)" -o ../$(BUILD-DEST)/make-torrc $(GO-PKG-PREFIX)/cmd/make-torrc

anondaemon: $(BUILD-DEST)
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) $(GO) build $(GO-BUILD-FLAGS) -ldflags="$(GO-LDFLAGS)" -o ../$(BUILD-DEST)/anondaemon $(GO-PKG-PREFIX)/cmd/anondaemon

safekill: $(BUILD-DEST)
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) $(GO) build $(GO-BUILD-FLAGS) -ldflags="$(GO-LDFLAGS)" -o ../$(BUILD-DEST)/safekill $(GO-PKG-PREFIX)/cmd/safekill

gui:
	$(CMAKE) -S gui -B $(GUI-BUILD-DIR) -DCMAKE_BUILD_TYPE=Release -DANONSURF_VERSION=$(VERSION)
	$(CMAKE) --build $(GUI-BUILD-DIR) --parallel

gui-check: gui
	$(CMAKE) --build $(GUI-BUILD-DIR) --target anonsurf-gui_qmllint

build: anonimax-cli dnstool make-torrc anondaemon safekill gui

check: gui-check
	# vendor/ is upstream code: only format-check our own sources.
	cd $(GO-SRC-DIR) && files=$$(find . -path ./vendor -prune -o -name '*.go' -print0 | xargs -0 gofmt -l) && \
		if [ -n "$$files" ]; then echo "not gofmt-clean: $$files"; exit 1; fi
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) $(GO) vet $(GO-PKG-PREFIX)/...
	cd $(GO-SRC-DIR) && $(GO-BUILD-ENV) GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GO-PKG-PREFIX)/cmd/...

inst-configs:
	cp configs/bridges.txt $(DESTDIR)/etc/anonimax/.
	cp configs/onion.pac $(DESTDIR)/etc/anonimax/.
	cp configs/anonimax.cfg $(DESTDIR)/etc/anonimax/.

	cp sys-units/animaxd.service $(DESTDIR)/lib/systemd/system/animaxd.service
	cp sys-units/anonimax-recovery.service $(DESTDIR)/lib/systemd/system/anonimax-recovery.service
	cp sys-units/anonimax-identity.service $(DESTDIR)/lib/systemd/system/anonimax-identity.service
	cp sys-units/tor@default.service.d/anonimax-recovery.conf $(DESTDIR)/lib/systemd/system/tor@default.service.d/anonimax-recovery.conf
	ln -sfn ../anonimax-recovery.service $(DESTDIR)/lib/systemd/system/multi-user.target.wants/anonimax-recovery.service

inst-bins:
	cp build/release/anonimax $(DESTDIR)/usr/bin/anonimax
	cp build/release/dnstool $(DESTDIR)/usr/bin/dnstool
	cp build/release/make-torrc $(DESTDIR)/usr/lib/anonimax/make-torrc
	cp build/release/anondaemon $(DESTDIR)/usr/lib/anonimax/anondaemon
	cp build/release/safekill $(DESTDIR)/usr/lib/anonimax/safekill

inst-gui:
	cp build/gui/anonsurf-gui $(DESTDIR)/usr/bin/anonimax-gui
	cp gui/org.parrotsec.Anonsurf.desktop $(DESTDIR)/usr/share/applications/.
	cp gui/org.parrotsec.AnonsurfCli.desktop $(DESTDIR)/usr/share/applications/.
	cp icons/anonsurf.png $(DESTDIR)/usr/share/icons/hicolor/512x512/apps/anonsurf.png

inst-perms:
	chmod +x $(DESTDIR)/usr/lib/anonimax/anondaemon $(DESTDIR)/usr/lib/anonimax/make-torrc $(DESTDIR)/usr/lib/anonimax/safekill $(DESTDIR)/usr/bin/anonimax $(DESTDIR)/usr/bin/anonimax-gui
	chmod 600 $(DESTDIR)/etc/anonimax/anonimax.cfg

mk-folders:
	mkdir -p $(DESTDIR)/etc/anonimax/
	mkdir -p $(DESTDIR)/usr/lib/anonimax/
	mkdir -p $(DESTDIR)/usr/bin/
	mkdir -p $(DESTDIR)/usr/share/applications/
	mkdir -p $(DESTDIR)/usr/share/icons/hicolor/512x512/apps/
	mkdir -p $(DESTDIR)/lib/systemd/system/
	mkdir -p $(DESTDIR)/lib/systemd/system/multi-user.target.wants/
	mkdir -p $(DESTDIR)/lib/systemd/system/tor@default.service.d/

install: mk-folders inst-bins inst-gui inst-configs inst-perms

uninstall:
	rm -rf /etc/anonimax/
	rm -rf /usr/lib/anonimax/
	rm /usr/bin/anonimax
	rm /usr/bin/anonimax-gui
	rm /usr/share/applications/org.parrotsec.Anonsurf.desktop
	rm /usr/share/applications/org.parrotsec.AnonsurfCli.desktop
	rm /usr/share/icons/hicolor/512x512/apps/anonsurf.png
	rm /lib/systemd/system/animaxd.service
	rm /lib/systemd/system/anonimax-recovery.service
	rm /lib/systemd/system/anonimax-identity.service
	rm /lib/systemd/system/multi-user.target.wants/anonimax-recovery.service
	rm /lib/systemd/system/tor@default.service.d/anonimax-recovery.conf
