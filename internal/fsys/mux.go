package fsys

import (
	"fmt"
	"io/fs"
	"sync"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
)

var (
	_instance *singleton
	lock      = &sync.Mutex{}
)

type singleton struct {
	mux fsimpl.FSMux
}

func (s *singleton) FS(uri string) (fs.FS, error) {
	fs, err := s.mux.Lookup(uri)
	if err != nil {
		return nil, fmt.Errorf("error obtaining the filesystem: '%w'", err)
	}

	return fs, nil
}

func getInstance() *singleton {
	if _instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if _instance == nil {
			_instance = &singleton{
				mux: fsimpl.FSMux{},
			}
			_instance.mux.Add(filefs.FS)
			_instance.mux.Add(httpfs.FS)
			_instance.mux.Add(gitfs.FS)
		}
	}

	return _instance
}
