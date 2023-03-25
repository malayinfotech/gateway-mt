ARG DOCKER_ARCH
FROM ${DOCKER_ARCH:-amd64}/alpine
ARG TAG
ARG GOARCH
ENV GOARCH ${GOARCH}
EXPOSE 20020
WORKDIR /app
VOLUME /root/.local/share/storx/linksharing

# mailcap adds mime type files we need for linksharing to determine mime types
# of files automatically.
RUN apk --no-cache --update add mailcap
COPY release/${TAG}/linksharing_linux_${GOARCH:-amd64} /app/linksharing
COPY pkg/linksharing/web/ /app/pkg/linksharing/web/
COPY cmd/linksharing/entrypoint /entrypoint
ENTRYPOINT ["/entrypoint"]
ENV STORX_CONFIG_DIR=/root/.local/share/storx/linksharing

# Healthcheck URL:
# http://<host>:20020/health/process
