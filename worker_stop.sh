#!/usr/bin/env bash
set -euo pipefail

# Gunakan kill -PGID agar make dan anak-anaknya ikut mati.
for f in /tmp/parser.pid /tmp/checker.pid /tmp/joiner.pid /tmp/executor.pid; do
  [ -f "$f" ] || continue
  pid=$(cat "$f")
  kill -- -"$pid" 2>/dev/null || true
  rm -f "$f"
done
