package dkvolume

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// DefaultDockerRootDirectory is the default directory where volumes will be created.
	DefaultDockerRootDirectory = "/var/lib/docker/volumes"

	defaultContentTypeV1          = "appplication/vnd.docker.plugins.v1+json"
	defaultImplementationManifest = `{"Implements": ["VolumeDriver"]}`
	pluginSpecDir                 = "/usr/share/docker/plugins"

	activatePath    = "/Plugin.Activate"
	createPath      = "/VolumeDriver.Create"
	remotePath      = "/VolumeDriver.Remove"
	hostVirtualPath = "/VolumeDriver.Path"
	mountPath       = "/VolumeDriver.Mount"
	unmountPath     = "/VolumeDriver.Unmount"
)

// Request is the structure that docker's requests are deserialized to.
type Request struct {
	Name string
}

// Response is the strucutre that the plugin's responses are serialized to.
type Response struct {
	Mountpoint string
	Err        string
}

// Driver represent the interface a driver must fulfill.
type Driver interface {
	Create(Request) Response
	Remove(Request) Response
	Path(Request) Response
	Mount(Request) Response
	Unmount(Request) Response
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	driver Driver
	mux    *http.ServeMux
}

type actionHandler func(Request) Response

// NewHandler initializes the request handler with a driver implementation.
func NewHandler(driver Driver) *Handler {
	h := &Handler{driver, http.NewServeMux()}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.mux.HandleFunc(activatePath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", defaultContentTypeV1)
		fmt.Fprintln(w, defaultImplementationManifest)
	})

	h.handle(createPath, func(req Request) Response {
		return h.driver.Create(req)
	})

	h.handle(remotePath, func(req Request) Response {
		return h.driver.Remove(req)
	})

	h.handle(hostVirtualPath, func(req Request) Response {
		return h.driver.Path(req)
	})

	h.handle(mountPath, func(req Request) Response {
		return h.driver.Mount(req)
	})

	h.handle(unmountPath, func(req Request) Response {
		return h.driver.Unmount(req)
	})
}

func (h *Handler) handle(name string, actionCall actionHandler) {
	h.mux.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(w, r)
		if err != nil {
			return
		}

		res := actionCall(req)

		encodeResponse(w, res)
	})
}

// ServeTCP makes the handler to listen for request in a given TCP address.
// It also writes the spec file on the right directory for docker to read.
func (h *Handler) ServeTCP(pluginName, addr string) error {
	return h.listenAndServe("tcp", addr, pluginName)
}

// ServeUnix makes the handler to listen for requests in a unix socket.
// It also creates the socket file on the right directory for docker to read.
func (h *Handler) ServeUnix(systemGroup, addr string) error {
	return h.listenAndServe("unix", addr, systemGroup)
}

func (h *Handler) listenAndServe(proto, addr, group string) error {
	server := http.Server{
		Addr:    addr,
		Handler: h.mux,
	}

	if err := os.MkdirAll(pluginSpecDir, 0755); err != nil {
		return err
	}

	start := make(chan struct{})

	var l net.Listener
	var err error
	switch proto {
	case "tcp":
		l, err = newTCPSocket(addr, nil, start)
		if err == nil {
			err = writeSpec(group, l.Addr().String())
		}
	case "unix":
		l, err = newUnixSocket(fullSocketAddr(addr), group, start)
	}
	if err != nil {
		return err
	}

	close(start)
	return server.Serve(l)
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (req Request, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func encodeResponse(w http.ResponseWriter, res Response) {
	w.Header().Set("Content-Type", defaultContentTypeV1)
	if res.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}

func writeSpec(name, addr string) error {
	spec := filepath.Join(pluginSpecDir, name+".spec")
	url := "tcp://" + addr
	return ioutil.WriteFile(spec, []byte(url), 0644)
}

func fullSocketAddr(addr string) string {
	if filepath.IsAbs(addr) {
		return addr
	}

	return filepath.Join(pluginSpecDir, addr+".sock")
}
