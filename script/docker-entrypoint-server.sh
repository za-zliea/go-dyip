#!/bin/bash

if [ -z ${CONFIG_DIR} ]; then
    CONFIG_DIR=/etc/dyip
fi

if [ ! -d ${CONFIG_DIR} ]; then
    mkdir ${CONFIG_DIR}
fi

if [[ $# = 1 ]] && [[ "$1" = 'dyip-server' ]]; then
    if [ ! -f ${CONFIG_DIR}/server.conf ]; then
        dyip-server -g -c ${CONFIG_DIR}/server.conf
    else
        dyip-server -c ${CONFIG_DIR}/server.conf
    fi
else
    exec "$@"
fi