# Docker volume nfs mounter

An example of the [Docker volume extension api](https://github.com/calavera/docker-volume-api)

Docker volume extension that NFS mounts a remote FS into your container

## Usage

`make build` to build the container

`make run` to run the volume plugin in a container, listening to the socket in the default
`/usr/share/docker/plugins/` dir.

To use the plugin when mounting an NFS export `nfs://127.0.0.1:/data`, run:

`docker run --rm -it --volume-driver=nfs -v 127.0.0.1/data:/no busybox ls -la`

> Note: because of the way docker parses colons, you need to omit them from the NFS share.

## Build and run in Boot2Docker qemu VM

`make iso` will build an experimental boot2docker ISO which will auto-start the nfs volume-plugin.

If you're on Linux, you can run it in Qemu-kvm using `make run`.

## License

MIT
