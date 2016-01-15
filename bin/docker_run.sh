#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
HOME_DIR="${SCRIPT_DIR}/.."
docker run "$@" -d -p 58080:8080 -v ${HOME_DIR}:/go/src/github.com/ugarcia/go_test_server go_test_server
