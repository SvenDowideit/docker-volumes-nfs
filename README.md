# Docker volume nfs mounter

An example of the [Docker volume extension api](https://github.com/calavera/docker-volume-api)

Docker volume extension that NFS mounts a remote FS into your container

## Usage

`make build` to build the container

`make run` to run the volume plugin in a container, listening to the socket in the default
`/usr/share/docker/plugins/` dir.

To use the plugin when mounting a volume, run:



## License

MIT
