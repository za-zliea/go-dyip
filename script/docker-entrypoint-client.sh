#!/bin/bash

if [ -z ${CONFIG_DIR} ]; then
    CONFIG_DIR=/etc/dyip
fi

if [ ! -d ${CONFIG_DIR} ]; then
    mkdir ${CONFIG_DIR}
fi

if [[ $# = 1 ]] && [[ "$1" = 'dyip-client' ]]; then
    if [ ! -f ${CONFIG_DIR}/client.conf ]; then
        dyip-client -g -c ${CONFIG_DIR}/client.conf
    else
        dyip-client -c ${CONFIG_DIR}/client.conf
    fi
else
    exec "$@"
fi