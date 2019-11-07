FROM alpine:3.10

RUN apk update ; \
	apk upgrade ; \
	rm -rf /var/cache/apk/* \
	export RANDOM_PASSWORD=`tr -dc A-Za-z0-9 < /dev/urandom | head -c71` ; \
	echo "root:$RANDOM_PASSWORD" | chpasswd ; \
	unset RANDOM_PASSWORD ; \
	passwd -l root

RUN set -ex; \
	apkArch="$(apk --print-arch)"; \
	case "$apkArch" in \
	armhf) arch='armv6' ;; \
	aarch64) arch='arm64' ;; \
	x86_64) arch='amd64' ;; \
	*) echo >&2 "error: unsupported architecture: $apkArch"; exit 1 ;; \
	esac; \
	wget --quiet -O /usr/local/bin/titanic "https://gitlab.com/hyperd/titanic/raw/master/releases/titanic-linux-$arch"; \
	chmod +x /usr/local/bin/titanic

RUN adduser -D -g '' appuser

# USER appuser

EXPOSE 3000 8443

CMD ["sh", "-c", "titanic --database.type=\"inmemory\""]

# Metadata
LABEL org.opencontainers.image.vendor="Hyperd" \
	org.opencontainers.image.url="https://hyperd.sh" \
	org.opencontainers.image.title="Titanic" \
	org.opencontainers.image.description="Container Solution API-exercise" \
	org.opencontainers.image.version="v0.5" \
	org.opencontainers.image.documentation="https://gitlab.com/hyperd/titanic/blob/master/README.md"

#sha256:29a159bc4f91abefd5cfcdae8236d3151c319925989034023973febec553af63