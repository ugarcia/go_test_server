#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
HOME_DIR="${SCRIPT_DIR}/.."
docker build "$@" -t go_test_server $HOME_DIR
