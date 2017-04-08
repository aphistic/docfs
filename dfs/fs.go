package dfs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	fusefs "bazil.org/fuse/fs"

	"github.com/aphistic/docfs/dfs/db"
	"github.com/efritz/glock"
)

type nodeType int

const (
	nRoot nodeType = iota
	nTags
	nTag
	nDocs
	nDoc
	nYear
	nMonth
	nDay
	nScratch
)

type DocFS struct {
	inodeLock  sync.Mutex
	inodeCount uint64
	inodes     map[nodeType]map[uint64]uint64

	fsRoot string
	root   *root
	fsdb   *db.DB

	clock glock.Clock
}

type node struct {
	fs *DocFS

	inode uint64
	name  string
}

func NewDocFS(fsRoot string) (*DocFS, error) {
	fInfo, err := os.Stat(fsRoot)
	if err == os.ErrNotExist || fInfo == nil {
		return nil, os.ErrNotExist
	} else if !fInfo.IsDir() {
		return nil, errors.New("Not a directory")
	}

	dbPath := path.Join(fsRoot, "docfs.db")
	fsdb, err := db.Open(dbPath)
	if err != nil {
		return nil, err
	}

	fs := &DocFS{
		inodes: make(map[nodeType]map[uint64]uint64),

		fsRoot: fsRoot,

		clock: glock.NewRealClock(),
	}
	fs.root = newRoot(fs)
	fs.fsdb = fsdb

	return fs, nil
}

func (f *DocFS) Root() (fusefs.Node, error) {
	return f.root, nil
}

func (f *DocFS) getInode(node nodeType, id uint64) uint64 {
	fmt.Printf("docfs: %v\n", f)
	f.inodeLock.Lock()
	defer f.inodeLock.Unlock()

	typeMap, ok := f.inodes[node]
	if !ok {
		typeMap = make(map[uint64]uint64)
		f.inodes[node] = typeMap
	}

	inode, ok := typeMap[id]
	if !ok {
		f.inodeCount++
		inode = f.inodeCount
		typeMap[id] = inode
	}

	fmt.Printf("Nodes: %+v\n", f.inodes)

	return inode
}

func (f *DocFS) openScratch(id uint64) (*os.File, error) {
	scratchRoot := path.Join(f.fsRoot, "scratch")
	docPath := path.Join(scratchRoot, fmt.Sprintf("%d", id))
	fmt.Printf("scratch path: %s\n", docPath)

	if _, err := os.Stat(scratchRoot); os.IsNotExist(err) {
		err = os.MkdirAll(scratchRoot, 0755)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(docPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *DocFS) Close() error {
	return f.fsdb.Close()
}
