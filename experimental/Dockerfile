FROM boot2docker/boot2docker

ENV LOCALBIN $ROOTFS/usr/local/bin/

# TODO is there an ADD attr that chmod's?
RUN curl -L -o $LOCALBIN/docker https://experimental.docker.com/builds/Linux/x86_64/docker-latest \
	&& chmod 755 $LOCALBIN/docker

# add the magic nfs volume driver
ADD ./docker-volume-nfs $LOCALBIN/docker-volume-nfs
# Create the plugins socket dir
RUN mkdir -p $ROOTFS/run/docker/plugins
# TODO: should test to see if there is one, and merge
ADD ./bootsync.sh $ROOTFS/var/lib/boot2docker/bootsync.sh

#ADD https://github.com/docker/machine/releases/download/v0.1.0-rc1/docker-machine_linux_amd64 $LOCALBIN/machine
#ADD https://github.com/docker/swarm/releases/download/v0.1.0-rc1/docker-swarm-Linux-x86_64 $LOCALBIN/swarm
#ADD https://github.com/docker/fig/releases/download/1.1.0-rc1/docker-compose-Linux-x86_64 $LOCALBIN/docker-compose

RUN /make_iso.sh

CMD ["cat", "boot2docker.iso"]
