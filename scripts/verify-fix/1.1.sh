#!/usr/bin/env bash
source "$(dirname "$0")/lib.sh"
run_fix_check "1.1" "1.1 SQL-инъекция — Исправление" "^TestRound1_SQLInjection$"
