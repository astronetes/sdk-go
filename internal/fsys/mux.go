package fsys

import (
	"io/fs"
	"sync"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
)

var (
	instance *singleton
	lock     = &sync.Mutex{}
)

type singleton struct {
	mux fsimpl.FSMux
}

func (s *singleton) FS(uri string) (fs.FS, error) {
	return s.mux.Lookup(uri)
}

func getInstance() *singleton {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &singleton{
				mux: fsimpl.FSMux{},
			}
			instance.mux.Add(filefs.FS)
			instance.mux.Add(httpfs.FS)
			instance.mux.Add(gitfs.FS)
		}
	}
	return instance
}
