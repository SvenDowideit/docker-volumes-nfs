package main

// TODO: check that we're running with mount priv

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/calavera/dkvolume"
)

var (
	root = flag.String("root", dkvolume.DefaultDockerRootDirectory, "Docker volumes root directory")
)

type garbageDriver struct {
	root string
}

func (g garbageDriver) Create(r dkvolume.Request) dkvolume.Response {
	fmt.Printf("Create %v\n", r)
	return dkvolume.Response{}
}

func (g garbageDriver) Remove(r dkvolume.Request) dkvolume.Response {
	fmt.Printf("Remove %v\n", r)
	return dkvolume.Response{}
}

func (g garbageDriver) Path(r dkvolume.Request) dkvolume.Response {
	fmt.Printf("Path %v\n", r)
	return dkvolume.Response{Mountpoint: filepath.Join(g.root, r.Name)}
}

func (g garbageDriver) Mount(r dkvolume.Request) dkvolume.Response {
	p := filepath.Join(g.root, r.Name)

	// Docker volumes syntax doesn't allow
	//  using ":" so we should specify nfs mounts
	//  as "host_or_ip/path/to/nfs". Here we
	//  rebuild the right address host_or_ip:/path/to/nfs
	v := strings.Split(r.Name, "/")
	v[0] = v[0] + ":"
	source := strings.Join(v, "/")

	fmt.Printf("Mount %s at %s\n", source, p)

	if err := os.MkdirAll(p, 0755); err != nil {
		return dkvolume.Response{Err: err.Error()}
	}

	//if err := ioutil.WriteFile(filepath.Join(p, "test"), []byte("TESTTEST"), 0644); err != nil {
	//fmt.Printf("wrote %s\n", filepath.Join(p, "test"))
	// if err := run("mount", "--bind", "/data/ISOs", p); err != nil {
	if err := run("mount", source, p); err != nil {
		return dkvolume.Response{Err: err.Error()}
	}

	return dkvolume.Response{Mountpoint: p}
}

func (g garbageDriver) Unmount(r dkvolume.Request) dkvolume.Response {
	p := filepath.Join(g.root, r.Name)
	fmt.Printf("Unmount %s\n", p)

	if err := run("umount", p); err != nil {
		return dkvolume.Response{Err: err.Error()}
	}

	err := os.RemoveAll(p)
	return dkvolume.Response{Err: err.Error()}
}

func main() {
	d := garbageDriver{*root}
	h := dkvolume.NewHandler(d)
	fmt.Println(h.ServeUnix("", "/usr/share/docker/plugins/nfs.sock"))
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
