

build:
	docker build -t docker-volumes-nfs .
	docker create --name test-vnfs docker-volumes-nfs
	docker cp test-vnfs:/go/bin/app .
	docker rm test-vnfs

run:
	sudo ./app

containerrun:
	docker run --rm -it --privileged \
		-v /usr/share/docker/plugins/:/usr/share/docker/plugins/ \
		-v /var/lib/docker/volumes/:/var/lib/docker/volumes \
		-v /data:/data \
		docker-volumes-nfs

test:
	docker run --rm -it --volume-driver=nfs -v Users/:/no busybox ls -la /no
	
nfs:
	docker run --rm -it --volume-driver=nfs -v 127.0.0.1/data:/no busybox ls -la
