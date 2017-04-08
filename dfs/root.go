package dfs

import (
	"os"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type root struct {
	node

	fs   *DocFS
	tags *fsTags
	docs *fsDocs
}

func newRoot(fs *DocFS) *root {
	r := &root{
		fs: fs,
	}
	r.inode = fs.getInode(nRoot, 0)
	r.name = "root"

	r.tags = newFsTags(fs)
	r.docs = newFsDocs(fs)

	return r
}

func (r *root) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = r.inode
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (r *root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var children []fuse.Dirent
	children = append(children, fuse.Dirent{
		Inode: r.tags.inode,
		Name:  "tags",
	})
	children = append(children, fuse.Dirent{
		Inode: r.docs.inode,
		Name:  "documents",
	})
	return children, nil
}

func (r *root) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	if name == "tags" {
		return r.tags, nil
	} else if name == "documents" {
		return r.docs, nil
	}
	return nil, fuse.ENOENT
}
