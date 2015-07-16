# Docker volume extension api.

Go handler to create external volume extensions for Docker.

## Usage

This library is designed to be integrated in your program.

1. Implement the `VolumeDriver` interface.
2. Initialize a `VolumeHander` with your implementation.
2. Call the method `ListenAndServe` from the `VolumeHandler`.

```go
  d := MyVolumeDriver{}
  h := volumeapi.NewVolumeHandler(d)
  h.ListenAndServe("tcp", ":8080", "")
```

See a full example in https://github.com/calavera/docker-volume-keywhiz-fs

## License

MIT
