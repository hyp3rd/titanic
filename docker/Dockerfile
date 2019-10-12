FROM alpine:3.10

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

USER appuser

EXPOSE 3000 8443

ENTRYPOINT ["titanic"]

# Metadata
LABEL org.opencontainers.image.vendor="Hyperd" \
	org.opencontainers.image.url="https://hyperd.sh" \
	org.opencontainers.image.title="Titanic" \
	org.opencontainers.image.description="Container Solution API-exercise" \
	org.opencontainers.image.version="v0.5" \
	org.opencontainers.image.documentation="https://gitlab.com/hyperd/titanic/blob/master/README.md"
