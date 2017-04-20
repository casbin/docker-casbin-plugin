.PHONY: all binary install clean uninstall

LIBDIR=${DESTDIR}/lib/systemd/system
BINDIR=${DESTDIR}/usr/lib/docker

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

uninstall:
	rm -f ${LIBDIR}/casbin-authz-plugin.service
	rm -f ${LIBDIR}/casbin-authz-plugin.socket
	rm -f ${BINDIR}/casbin-authz-plugin
	rm -f ${BINDIR}/casbin.conf
	rm -f ${BINDIR}/basic_model.conf
	rm -f ${BINDIR}/basic_policy.csv
