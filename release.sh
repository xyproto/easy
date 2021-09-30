#!/bin/sh
#
# Create release tarballs/zip for 64-bit linux, BSD and Plan9 + 64-bit ARM + raspberry pi 2/3
#
name=easy
version=$(./version.sh -s)
echo "Version $version"
echo 'Compiling...'
export GOARCH=amd64
echo '* Linux'
GOOS=linux go build -mod=vendor -o $name.linux
echo '* Linux ARM64'
GOOS=linux GOARCH=arm64 go build -mod=vendor -o $name.linux_arm64
echo '* RPI 2/3/4'
GOOS=linux GOARCH=arm GOARM=7 go build -mod=vendor -o $name.rpi
echo '* Linux static w/ upx'
CGO_ENABLED=0 GOOS=linux go build -mod=vendor -v -trimpath -ldflags "-s" -a -o $name.linux_static && upx $name.linux_static

# Compress the Linux releases with xz
for p in linux linux_arm64 rpi linux_static; do
  echo "Compressing $name-$version.$p.tar.xz"
  mkdir "$name-$version-$p"
  cp $name.1 "$name-$version-$p/"
  gzip "$name-$version-$p/$name.1"
  cp $name.$p "$name-$version-$p/$name"
  cp COPYING "$name-$version-$p/"
  tar Jcf "$name-$version-$p.tar.xz" "$name-$version-$p/"
  rm -r "$name-$version-$p"
  rm $name.$p
done
