#!/bin/bash
set -e

if [ -f .env ]; then
    source .env
fi

go build -o ./build/voltcraft src/*.go
sudo -E ./build/voltcraft

