

build:
	docker build -t docker-volumes-nfs .
	docker create --name test-vnfs docker-volumes-nfs
	docker cp test-vnfs:/go/bin/app .
	docker rm test-vnfs

run:
	sudo ./app

containerrun:
	docker run --rm -it --privileged \
		-v /run/docker/plugins:/run/docker/plugins \
		-v /var/lib/docker/volumes/:/var/lib/docker/volumes \
		-v /data:/data \
		docker-volumes-nfs

test:
	docker run --rm -it --volume-driver=nfs -v Users/:/no busybox ls -la /no
	
nfs:
	docker run --rm -it --volume-driver=nfs -v 127.0.0.1/data:/no busybox ls -la

iso: build
	cp app experimental/docker-volume-nfs
	docker build -t boot2docker:experimental experimental/
	docker run --name boot2docker-run boot2docker:experimental > experimental/boot2docker.iso
	docker cp boot2docker-run:/linux-kernel/arch/x86_64/boot/bzImage experimental/
	docker cp boot2docker-run:/tmp/iso/boot/vmlinuz64 experimental/
	docker cp boot2docker-run:/tmp/iso/boot/initrd.img experimental/
	docker rm boot2docker-run

run: 
	qemu-system-x86_64 -serial stdio \
		-curses \
		-net nic,vlan=0 -net user,vlan=0 \
		-m 2048M \
		-kernel experimental/vmlinuz64 -initrd experimental/initrd.img \
		-append "root=/dev/ram0 rw sven=test panic=0 append loglevel=7 user=docker console=ttyAMA0 console=ttyS0 console=tty0 apparmor=0 selinux=1 noembed nomodeset norestore base"

