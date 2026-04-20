#!/usr/bin/env bash
source "$(dirname "$0")/lib.sh"
run_fix_check "1.4" "1.4 Обход JWT-аутентификации — Исправление" "^TestRound1_JWTValidation$"
