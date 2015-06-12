package main

// TODO: check that we're running with mount priv

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/calavera/docker-volume-api"
)

var (
	root = flag.String("root", volumeapi.DefaultDockerRootDirectory, "Docker volumes root directory")
)

type garbageDriver struct {
	root string
}

func (g garbageDriver) Create(r volumeapi.VolumeRequest) volumeapi.VolumeResponse {
	fmt.Printf("Create %v\n", r)
	return volumeapi.VolumeResponse{}
}

func (g garbageDriver) Remove(r volumeapi.VolumeRequest) volumeapi.VolumeResponse {
	fmt.Printf("Remove %v\n", r)
	return volumeapi.VolumeResponse{}
}

func (g garbageDriver) Path(r volumeapi.VolumeRequest) volumeapi.VolumeResponse {
	fmt.Printf("Path %v\n", r)
	return volumeapi.VolumeResponse{Mountpoint: filepath.Join(g.root, r.Name)}
}

func (g garbageDriver) Mount(r volumeapi.VolumeRequest) volumeapi.VolumeResponse {
	p := filepath.Join(g.root, r.Name)

	v := strings.Split(r.Name, "/")
	v[0] = v[0]+":"
	source := strings.Join(v, "/")

	fmt.Printf("Mount %s at %s\n", source, p)

	if err := os.MkdirAll(p, 0755); err != nil {
		return volumeapi.VolumeResponse{Err: err}
	}

	//if err := ioutil.WriteFile(filepath.Join(p, "test"), []byte("TESTTEST"), 0644); err != nil {
	//fmt.Printf("wrote %s\n", filepath.Join(p, "test"))
	// if err := run("mount", "--bind", "/data/ISOs", p); err != nil {
	if err := run("mount", source, p); err != nil {
		return volumeapi.VolumeResponse{Err: err}
	}

	return volumeapi.VolumeResponse{Mountpoint: p}
}

func (g garbageDriver) Unmount(r volumeapi.VolumeRequest) volumeapi.VolumeResponse {
	p := filepath.Join(g.root, r.Name)
	fmt.Printf("Unmount %s\n", p)

	if err := run("umount", p); err != nil {
		return volumeapi.VolumeResponse{Err: err}
	}

	err := os.RemoveAll(p)
	return volumeapi.VolumeResponse{Err: err}
}

func main() {
	d := garbageDriver{*root}
	h := volumeapi.NewVolumeHandler(d)
	fmt.Println(h.ListenAndServe("unix", "/usr/share/docker/plugins/nfs.sock", ""))
}

var (
	verbose = true
)

func run(exe string, args ...string) error {
        cmd := exec.Command(exe, args...)
        if verbose {
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr
                fmt.Printf("executing: %v %v", exe, strings.Join(args, " "))
        }
        if err := cmd.Run(); err != nil {
                return err
        }
        return nil
}

