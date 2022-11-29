#!/bin/bash

set -e

if [ $# -eq 0 ]; then
	echo "Usage: ./update.sh v#.#.#"
	exit
fi

VERSION=$1

docker build -t sfs-builder -f ./Dockerfile.all .

ID=$(docker create sfs-builder)

rm -rf out
mkdir -p out
docker cp "${ID}:/build/pkg/linux-amd64/serve" "./out/static-file-server-${VERSION}-linux-amd64"
docker cp "${ID}:/build/pkg/linux-i386/serve" "./out/static-file-server-${VERSION}-linux-386"
docker cp "${ID}:/build/pkg/linux-arm6/serve" "./out/static-file-server-${VERSION}-linux-arm6"
docker cp "${ID}:/build/pkg/linux-arm7/serve" "./out/static-file-server-${VERSION}-linux-arm7"
docker cp "${ID}:/build/pkg/linux-arm64/serve" "./out/static-file-server-${VERSION}-linux-arm64"
docker cp "${ID}:/build/pkg/darwin-amd64/serve" "./out/static-file-server-${VERSION}-darwin-amd64"
docker cp "${ID}:/build/pkg/win-amd64/serve.exe" "./out/static-file-server-${VERSION}-windows-amd64.exe"

docker rm -f "${ID}"
docker rmi sfs-builder

docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag "halverneus/static-file-server:${VERSION}" .
docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag halverneus/static-file-server:latest .

echo "Done"
