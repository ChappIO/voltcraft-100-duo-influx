#!/bin/bash
set -e

go build -o ./build/voltcraft src/*.go
sudo -E ./build/voltcraft

