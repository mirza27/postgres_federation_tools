#!/usr/bin/env bash
set -euo pipefail

mkdir -p logs

# Jalankan tiap worker di process group sendiri supaya mudah dibunuh dengan kill -PGID
setsid bash -c "exec make parser   > logs/parser.log  2>&1" &
echo $! > /tmp/parser.pid

setsid bash -c "exec make joiner   > logs/joiner.log  2>&1" &
echo $! > /tmp/joiner.pid

setsid bash -c "exec make checker  > logs/checker.log 2>&1" &
echo $! > /tmp/checker.pid

setsid bash -c "exec make executor > logs/executor.log 2>&1" &
echo $! > /tmp/executor.pid
