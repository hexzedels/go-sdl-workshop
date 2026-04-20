#!/usr/bin/env bash
source "$(dirname "$0")/lib.sh"
run_fix_check "2.2" "2.2 Слабый генератор случайных чисел — Исправление" "^TestRound2_WeakRandom$"
