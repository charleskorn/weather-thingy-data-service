#!/usr/bin/env bash

set -e

BUILDDIR=./rkt_layout
MANIFEST=./rkt/manifest.json
OUTPUTFILE=weather-thingy-data-service.aci

echo = Validating $MANIFEST...
actool -debug validate $MANIFEST

echo = Creating directory structure...

rm -rf $BUILDDIR
mkdir -p $BUILDDIR
mkdir -p $BUILDDIR/rootfs
mkdir -p $BUILDDIR/rootfs/bin

echo = Copying files into place...
cp $MANIFEST $BUILDDIR/manifest

echo = Building data service...
CGO_ENABLED=0 GOOS=linux go build -o $BUILDDIR/rootfs/bin/weather-thingy-data-service -a -installsuffix cgo .

echo = Building image...
actool build --overwrite $BUILDDIR $OUTPUTFILE

echo = Validating image...
actool -debug validate $OUTPUTFILE

echo = Done. Image $OUTPUTFILE created.
