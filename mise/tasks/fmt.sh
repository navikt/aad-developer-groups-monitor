#!/usr/bin/env bash
#MISE description="Format go files using gofumpt"
#MISE wait_for=["gofix"]
set -euo pipefail

go tool mvdan.cc/gofumpt -w ./
