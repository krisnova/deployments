package thenovadiary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/kris-nova/logger"
)

// TODO We need to support native times in our records
// TODO  We need this to work well with unit tests for our crons

const (
	DatabaseMode = 0755
	DatabaseName = "cache.db"
	DatabaseDir  = "thenovadiary"
)

type Cache struct {
	Name        string
	path        *os.File
	Records     map[string]*Record `json:"Records"`
	globalMutex sync.Mutex
}

type Record struct {
	Found bool
	Key   string
	Value interface{}
}

func NewCache(name string) *Cache {

	// System to locate the cache on disk
	path, err := GetCachePath(name)
	if err != nil {
		logger.Warning("unable to find cache path using non deterministic path: %s", err)
		// Assume linux
		file, err := ioutil.TempFile("/tmp", "thenovadiary_")
		if err != nil {
			logger.Critical("critical error: unable to find cache persistent store: %v", err)
			return nil
		}
		path = file
	}
	return &Cache{
		path:    path,
		Name:    name,
		Records: make(map[string]*Record),
	}
}

var cachedPath *os.File

// GetCachePath is a deterministic function
// to return an *os.File based on the system
// running.
//
// By design this will also ensure the path is
// writeable and created by this process's file
// descriptor access.
func GetCachePath(name string) (*os.File, error) {
	if cachedPath != nil {
		return cachedPath, nil
	}
	getPath := func() string {
		// GOOS/GOARCH
		// https://github.com/golang/go/blob/master/src/go/build/syslist.go
		// const goosList = "aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris windows zos "
		// const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm "
		if runtime.GOOS == "linux" {
			logger.Info("Linux OS detected")
			home := os.Getenv("HOME")
			if home != "" {
				return path.Join(home, ".config", DatabaseDir, name, DatabaseName)
			} else {
				return path.Join("/var/lib/", DatabaseDir, name, DatabaseName)
			}
		} else if runtime.GOOS == "freebsd" {
			logger.Info("FreeBSD detected")
			// /var/lib/thenovadiary/cache.db
			return "/var/lib/thenovadiary/cache.db"

		}
		gopath := os.Getenv("GOPATH")
		if gopath != "" {
			return path.Join(gopath, ".config", DatabaseDir, name, DatabaseName)
		}
		// Default case (No Path Found!)
		return ""
	}
	path := getPath()
	if path == "" {
		return nil, fmt.Errorf("unable to find path")
	}

	// Ensure the path exists (dirname)
	pathDir := filepath.Dir(path)
	logger.Debug("Ensuring path: %s", pathDir)
	err := os.MkdirAll(pathDir, DatabaseMode)
	if err != nil {
		logger.Debug("unable to mkdir -p cache directory: %v", err)
	}

	_, err = os.Stat(path)
	// touch
	if os.IsNotExist(err) {
		// Database does not exist
		createdFile, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("unable to create new path %s: %v", path, err)
		}
		createdFile.Close()
	} else if err != nil {
		return nil, fmt.Errorf("unable to stat() %s: %v", path, err)
	}
	// Regardless of touch, open the file (sanity check)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to open path: %v", err)
	}
	cachedPath = f
	return f, nil
}

func (c *Cache) Clean() error {
	//logger.Info(c.Now())
	c.globalMutex.Lock()
	return os.Remove(c.path.Name())
}

func (c *Cache) Recover() (int, error) {
	//logger.Info(c.Now())
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()

	jBytes, err := ioutil.ReadFile(c.path.Name())
	if err != nil {
		return -1, fmt.Errorf("unable to read file %s: %v", c.path.Name(), err)
	}
	newC := Cache{}
	err = json.Unmarshal(jBytes, &newC)
	if err != nil {
		return -1, fmt.Errorf("unable JSON unmarshal file %s: %v", c.path.Name(), err)
	}

	// Manually rebuild parts of the cache
	delta := float64(len(newC.Records) - len(c.Records))
	if delta != float64(0) {
		// Logic to manage delta here
	}
	c.Records = newC.Records
	return int(delta), nil
}

func (c *Cache) Persist() error {
	//logger.Info(c.Now())
	jBytes, err := json.Marshal(&c)
	if err != nil {
		return fmt.Errorf("unable to persist (JSON) to disk: %v", err)
	}
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	err = ioutil.WriteFile(c.path.Name(), jBytes, DatabaseMode)
	if err != nil {
		return fmt.Errorf("unable to persist (WRITE) to disk: %v", err)
	}
	return nil
}

func (c *Cache) Get(key string) *Record {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	r := &Record{
		Found: false,
	}
	if value, ok := c.Records[key]; ok {
		value.Found = true
		return value
	}
	return r
}

func (c *Cache) Set(key string, r *Record) {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	c.Records[key] = r
}

func (c *Cache) Remove(key string) {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	if _, ok := c.Records[key]; ok {
		delete(c.Records, key)
	}
}

func (c *Cache) Now() string {
	c.globalMutex.Lock()
	defer c.globalMutex.Unlock()
	return time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
}
