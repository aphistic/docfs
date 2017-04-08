package dfs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	fusefs "bazil.org/fuse/fs"

	"github.com/aphistic/docfs/dfs/db"
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

	root *root
	fsdb *db.DB
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

func (f *DocFS) Close() error {
	return f.fsdb.Close()
}
