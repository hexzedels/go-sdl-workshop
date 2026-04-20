#!/usr/bin/env bash
source "$(dirname "$0")/lib.sh"
run_fix_check "3.3" "3.3 Атака по времени — Исправление" "^TestRound3_TimingAttack$"
