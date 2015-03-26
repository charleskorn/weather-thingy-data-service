#!/usr/bin/env bash

set -e

go-bindata -pkg main -o bindata.go db/migrations/
