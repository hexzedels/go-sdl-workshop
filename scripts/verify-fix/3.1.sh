#!/usr/bin/env bash
source "$(dirname "$0")/lib.sh"
run_fix_check "3.1" "3.1 Утечка горутин — Исправление" "^TestRound3_GoroutineLeak$"
