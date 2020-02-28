#!/bin/bash -ex

rm -rf bin && SKIP_VALIDATE=true ./scripts/ci
#GO111MODULE=on
#STATIC="-extldflags '-static'"
#
#LDFLAGS="
#    -X $PKG/pkg/version.Version=$VERSION
#    -X $PKG/pkg/version.GitCommit=${COMMIT:0:8}
#    -X $PKG/vendor/$PKG_CONTAINERD/version.Version=$VERSION_CONTAINERD
#    -X $PKG/vendor/$PKG_CONTAINERD/version.Package=$PKG_RANCHER_CONTAINERD
#    -X $PKG/vendor/$PKG_CRICTL/pkg/version.Version=$VERSION_CRICTL
#    -w -s
#"
#
#CGO_ENABLED=1 go build -ldflags "$LDFLAGS $STATIC" -o k3s 
#scp dist/artifacts/k3s root@destiny:/usr/local/bin
#scp dist/artifacts/k3s root@dream:/usr/local/bin
#scp dist/artifacts/k3s root@death:/usr/local/bin
#
