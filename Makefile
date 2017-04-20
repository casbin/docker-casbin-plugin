.PHONY: all binary install clean

LIBDIR=${DESTDIR}/lib/systemd/system
BINDIR=${DESTDIR}/usr/lib/docker/

all: binary

binary:
	go build  -o casbin-authz-plugin .

install:
	install -m 644 systemd/casbin-authz-plugin.service ${LIBDIR}
	install -m 644 systemd/casbin-authz-plugin.socket ${LIBDIR}
	install -m 755 casbin-authz-plugin ${BINDIR}
	install -m 644 casbin.conf ${BINDIR}
	install -m 644 examples/basic_model.conf ${BINDIR}
	install -m 644 examples/basic_policy.csv ${BINDIR}

clean:
	rm -f casbin-authz-plugin
