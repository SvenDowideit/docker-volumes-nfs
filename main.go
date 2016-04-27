package main

// TODO: check that we're running with mount priv

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	//"github.com/calavera/docker-volume-api"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	pluginId = "nfs"
)

var (
	socketDir = "/run/docker/plugins/"
	socketAddress = filepath.Join(socketDir, strings.Join([]string{pluginId, ".sock"}, ""))
	defaultDir    = filepath.Join(volume.DefaultDockerRootDirectory, strings.Join([]string{"_", pluginId}, ""))
	root          = flag.String("root", defaultDir, "NFS volumes root directory")
)

type nfsDriver struct {
	root string
}

func (g nfsDriver) Create(r volume.Request) volume.Response {
	fmt.Printf("Create %v\n", r)
	return volume.Response{}
}

func (g nfsDriver) Remove(r volume.Request) volume.Response {
	fmt.Printf("Remove %v\n", r)
	return volume.Response{}
}

func (g nfsDriver) Path(r volume.Request) volume.Response {
	fmt.Printf("Path %v\n", r)
	return volume.Response{Mountpoint: filepath.Join(g.root, r.Name)}
}

func (g nfsDriver) Get(r volume.Request) volume.Response {
	fmt.Printf("Get %v\n", r)
	return volume.Response{}
}

func (g nfsDriver) List(r volume.Request) volume.Response {
	fmt.Printf("List %v\n", r)
	return volume.Response{}
}

func (g nfsDriver) Mount(r volume.Request) volume.Response {
	p := filepath.Join(g.root, r.Name)

	v := strings.Split(r.Name, "/")
	v[0] = v[0] + ":"
	source := strings.Join(v, "/")

	fmt.Printf("Mount %s at %s\n", source, p)

	if err := os.MkdirAll(p, 0755); err != nil {
		return volume.Response{Err: err.Error()}
	}

	// if err := ioutil.WriteFile(filepath.Join(p, "test"), []byte("TESTTEST"), 0644); err != nil {
	// fmt.Printf("wrote %s\n", filepath.Join(p, "test"))
	// if err := run("mount", "--bind", "/data/ISOs", p); err != nil {
	if err := run("mount", "-o", "port=2049,nolock,proto=tcp", source, p); err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{Mountpoint: p}
}

func (g nfsDriver) Unmount(r volume.Request) volume.Response {
	p := filepath.Join(g.root, r.Name)
	fmt.Printf("Unmount %s\n", p)

	if err := run("umount", p); err != nil {
		return volume.Response{Err: err.Error()}
	}

	err := os.RemoveAll(p)
	return volume.Response{Err: err.Error()}
}

func main() {
	d := nfsDriver{*root}
	h := volume.NewHandler(d)
	fmt.Printf("listening on %s\n", socketAddress)
	fmt.Println(h.ServeUnix("root", socketAddress))
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
