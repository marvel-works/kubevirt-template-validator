#!/bin/sh

set -ex

if [ -z "$1" ]; then
	echo "usage: $0 <tag>"
	exit 1
fi

TAG="$1"  #TODO: validate tag is vX.Y.Z
COMPONENT="kubevirt-template-validator"
LDFLAGS="\
-X 'github.com/fromanirh/kubevirt-template-validator/internal/pkg/version.COMPONENT=$COMPONENT' \
-X 'github.com/fromanirh/kubevirt-template-validator/internal/pkg/version.VERSION=$TAG' "
if git rev-parse &>/dev/null; then
    BRANCH=$( git rev-parse --abbrev-ref HEAD )
    REVISION=$( git rev-parse --short HEAD )
    LDFLAGS="${LDFLAGS}\
-X 'github.com/fromanirh/kubevirt-template-validator/internal/pkg/version.BRANCH=$BRANCH' \
-X 'github.com/fromanirh/kubevirt-template-validator/internal/pkg/version.REVISION=$REVISION' "
fi

export GO111MODULE=on
export GOPROXY=off
export GOFLAGS=-mod=vendor
cd cmd/kubevirt-template-validator && \
    go build -v -ldflags="$LDFLAGS" .
