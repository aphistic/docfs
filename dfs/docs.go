package dfs

import (
	"fmt"
	"os"

	"strconv"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type fsDocs struct {
	node

	fs *DocFS
}

func newFsDocs(fs *DocFS) *fsDocs {
	t := &fsDocs{
		fs: fs,
	}
	t.inode = fs.getInode(nDocs, 0)
	t.name = "docs"

	return t
}

func (d *fsDocs) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = d.inode
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (d *fsDocs) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var children []fuse.Dirent

	years, err := d.fs.fsdb.GetYears()
	if err != nil {
		return nil, err
	}

	for _, year := range years {
		children = append(children, fuse.Dirent{
			Name:  fmt.Sprintf("%d", year),
			Inode: d.fs.getInode(nYear, uint64(year)),
		})
	}

	return children, nil
}

func (d *fsDocs) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fusefs.Node, error) {
	year, err := strconv.ParseUint(req.Name, 10, 0)
	if err != nil {
		return nil, fuse.EPERM
	}

	err = d.fs.fsdb.AddYear(year)
	if err != nil {
		return nil, err
	}

	return newFsYear(d.fs, year), nil
}

func (d *fsDocs) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	if !req.Dir {
		return fuse.EIO
	}

	year, err := strconv.ParseUint(req.Name, 10, 0)
	if err != nil {
		return fuse.EIO
	}

	err = d.fs.fsdb.RemoveYear(year)
	if err != nil {
		return err
	}

	return nil
}

func (d *fsDocs) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	year, err := strconv.ParseUint(name, 10, 0)
	if err != nil {
		return nil, fuse.ENOENT
	}

	year, err = d.fs.fsdb.GetYear(year)
	if err != nil {
		return nil, fuse.ENOENT
	}

	return newFsYear(d.fs, year), nil
}

type fsDoc struct {
	node

	fs *DocFS

	ID uint64
}

func newFsDoc(fs *DocFS, id uint64, name string) *fsDoc {
	d := &fsDoc{
		fs: fs,
		ID: id,
	}
	d.inode = fs.getInode(nDoc, id)
	d.name = name

	return d
}

func (d *fsDoc) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = d.inode
	attr.Mode = 0755
	return nil
}
