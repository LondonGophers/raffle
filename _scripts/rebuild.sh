#!/usr/bin/env bash

set -euo pipefail

go run github.com/gopherjs/gopherjs build -m -o docs/raffle.js
cp index.html docs
