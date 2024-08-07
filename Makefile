.POSIX:

PREFIX=/usr/
DESTDIR=/
BINARY_NAME=vib

all: build # plugins

build:
	mkdir -p build
	sed 's|$$INSTALLPREFIX$$|${PREFIX}|g' core/plugins.in > core/plugins.go
	go build -a -o build/${BINARY_NAME}

build-plugins: FORCE
	mkdir -p build/plugins
	$(MAKE) -C plugins/

install:
	install -Dm755 -t ${DESTDIR}/${PREFIX}/bin/ ./build/${BINARY_NAME}

install-plugins:
	install -Dm644 -t ${DESTDIR}/${PREFIX}/share/vib/plugins/ ./build/plugins/*.so

clean:
	rm -r build
	rm core/plugins.go

FORCE:
