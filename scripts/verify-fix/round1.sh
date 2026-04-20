#!/usr/bin/env bash
# round1.sh — back-compat wrapper that runs all Round 1 per-challenge scripts.
set -e
D="$(dirname "$0")"
bash "$D/1.1.sh"
bash "$D/1.3.sh"
bash "$D/1.4.sh"
