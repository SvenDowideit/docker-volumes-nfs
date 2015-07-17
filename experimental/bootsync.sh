#!/bin/sh

echo "start NFS Docker volume plugin"
/usr/local/bin/docker-volume-nfs > /var/log/docker-volume-nfs.log 2>&1 &
