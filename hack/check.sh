#!/bin/bash

set -e

check_fmt() {
  for pkg in `ls -d src/* | grep -v vendor | grep -v test`; do
    echo check $pkg
    check_gofmt_in $pkg
  done
}

check_gofmt_in() {
  files=`gofmt -l $1`
  if [ -n "$files" ]; then
    echo "$files" >&2
    exit 1
  fi
}

check_lint() {
  exec gometalinter \
    --fast \
    --disable=gotype \
    --vendor \
    --deadline=360s \
    --severity=golint:error \
    --errors \
    --skip=test \
    --exclude=_mock \
    --exclude=tests \
    ./...
}

case "$1" in
  fmt) check_fmt ;;
  lint) check_lint ;;
esac
