#!/bin/sh
#
# Script doing IO-related busywork.
# For testing.
#
while true; do
  mkdir -p _testdir
  sync
  echo asdf > _testdir/fil
  sync
  sed -i 's/asdf/qwerty/g' _testdir/fil
  sync
  rm -rf _testdir
  sync
done
