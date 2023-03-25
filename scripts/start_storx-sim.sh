#!/usr/bin/env bash
set -euo pipefail

go install -race -v \
    storx/cmd/satellite@latest \
    storx/cmd/storagenode@latest \
    storx/cmd/storx-sim@latest \
    storx/cmd/versioncontrol@latest \
    storx/cmd/uplink@latest \
    storx/cmd/identity@latest \
    storx/cmd/certificates@latest \
    storx/cmd/multinode@latest

go install -race -v gateway@latest

echo "=== Finished installing dependencies"

echo "=== Setting up storx-sim..."

until storx-sim -x --host sim network setup; do
    echo "*** redis/postgres is not yet available; waiting for 3s..."
    sleep 3
done

sed -i 's/# metainfo.multiple-versions: false/metainfo.multiple-versions: true/g' "$(storx-sim network env SATELLITE_0_DIR)/config.yaml"

echo "=== Enabled multiple versions"

sed -i 's/# metainfo.rate-limiter.enabled: true/metainfo.rate-limiter.enabled: false/g' "$(storx-sim network env SATELLITE_0_DIR)/config.yaml"

echo "=== Disabled rate limiting"

echo "=== Starting storx-sim..."

storx-sim -x network run
