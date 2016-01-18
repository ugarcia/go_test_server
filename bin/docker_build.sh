#!/usr/bin/env bash
IMAGE_NAME=go_test_server
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
HOME_DIR="${SCRIPT_DIR}/.."
cp ~/.ssh/id_rsa ${HOME_DIR}/id_rsa
docker build "$@" -t $IMAGE_NAME $HOME_DIR
rm -rf ${HOME_DIR}/id_rsa
