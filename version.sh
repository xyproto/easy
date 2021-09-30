#!/bin/sh -e
#
# Self-modifying script that updates the version numbers
#

# The current version goes here, as the default value
CURRENT_VERSION='1.4.0'
VERSION=${1:-$CURRENT_VERSION}

if [ "$1" == '-s' ]; then
  echo "$CURRENT_VERSION"
  exit 0
elif [ -z "$1" ]; then
  echo "The current version is $VERSION, pass the new version as the first argument if you wish to change it"
  exit 0
fi

echo "Setting the version to $VERSION"

# Set the version in various files
setconf README.md '* Version' $VERSION
setconf main.go versionString "\"easy "$VERSION"\""

# Update the version in the man page
tac easy.1 | sed "0,/[[:digit:]]*\.[[:digit:]]*\.[[:digit:]]*/s//$VERSION/" | tac > /tmp/easy.1 && mv -f -v /tmp/easy.1 easy.1

# Update the date in the man page
d=$(LC_ALL=C date +'%d %b %Y')
sed -i "s/\"[0-9]* [A-Z][a-z]* [0-9]*\"/\"$d\"/g" easy.1

# Update the version in this script
sed -i "s/[[:digit:]]*\.[[:digit:]]*\.[[:digit:]]*/$VERSION/g" "$0"
