VERSION 0.6
FROM golang:1.18
WORKDIR /go/storx-up

gatewaymt-bin:
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        go build -o release/earthly/gateway-mt storx.io/gateway-mt/cmd/gateway-mt
    SAVE ARTIFACT release/earthly binaries AS LOCAL release/earthly

authservice-bin:
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        go build -o release/earthly/authservice storx.io/gateway-mt/cmd/authservice
    SAVE ARTIFACT release/earthly binaries AS LOCAL release/earthly

linksharing-bin:
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        go build -o release/earthly/linksharing storx.io/gateway-mt/cmd/linksharing
    SAVE ARTIFACT release/earthly binaries AS LOCAL release/earthly

build-storxup:
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        CGO_ENABLED=0 go install storx.io/storx-up@main
    SAVE ARTIFACT /go/bin binaries AS LOCAL dist/up

deploy-remote:
    FROM ubuntu
    RUN apt-get update && apt-get install -y git wget unzip
    RUN cd /tmp && wget https://releases.hashicorp.com/nomad/1.3.5/nomad_1.3.5_linux_amd64.zip -O nomad.zip && unzip nomad.zip && mv nomad /usr/local/bin && rm nomad.zip
    COPY +build-storxup/binaries  /usr/local/bin
    COPY .git .git
    ARG TAG=$(git rev-parse --short HEAD)
    ARG IMAGE=img.dev.storx.io/dev/edge
    BUILD +build-tagged-image --TAG=$TAG
    ARG --required nomad
    ARG --required ip
    ENV NOMAD_ADDR=$nomad
    RUN storx-up init nomad --name=edge --ip=$ip edge
    RUN storx-up image gateway-mt,linksharing,authservice $IMAGE:$TAG
    RUN --push nomad run --verbose storx.hcl

build-image:
    FROM storxlabs/ci
    COPY .git .git
    ARG TAG=$(git rev-parse --short HEAD)
    BUILD +build-tagged-image --TAG=$TAG

build-tagged-image:
    ARG --required TAG
    ARG --required IMAGE
    FROM img.dev.storx.io/storxup/base:20221011-2
    COPY +gatewaymt-bin/binaries /var/lib/storx/go/bin/
    COPY +linksharing-bin/binaries /var/lib/storx/go/bin/
    COPY +authservice-bin/binaries /var/lib/storx/go/bin/
    COPY +build-storxup/binaries  /var/lib/storx/go/bin/
    COPY pkg/linksharing/web /var/lib/storx/pkg/linksharing/web
    SAVE IMAGE --push $IMAGE:$TAG $IMAGE:latest

run:
    LOCALLY
    RUN docker-compose up

check-format:
   COPY . .
   RUN mkdir build
   RUN bash -c '[[ $(git status --short) == "" ]] || (echo "Before formatting, please commit all your work!!! (Formatter will format only last commit)" && exit -1)'
   RUN git show --name-only --pretty=format: | grep ".go" | xargs --no-run-if-empty -n1 gofmt -s -w
   RUN git diff > build/format.patch
   SAVE ARTIFACT build/format.patch

format:
   LOCALLY
   COPY +check-format/format.patch build/format.patch
   RUN git apply --allow-empty build/format.patch
   RUN git status
