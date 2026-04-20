#!/usr/bin/env bash
source "$(dirname "$0")/lib.sh"
run_fix_check "4.1" "4.1 Секреты в образе — Исправление" "^TestRound4_HardcodedSecrets$"
