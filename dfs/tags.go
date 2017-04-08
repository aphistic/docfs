package dfs

import (
	"os"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type fsTags struct {
	node

	fs *DocFS
}

func newFsTags(fs *DocFS) *fsTags {
	t := &fsTags{
		fs: fs,
	}
	t.inode = fs.getInode(nTags, 0)
	t.name = "tags"

	return t
}

func (t *fsTags) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = t.inode
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (t *fsTags) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	tags, err := t.fs.fsdb.GetTags()
	if err != nil {
		return nil, err
	}

	var children []fuse.Dirent
	for _, tag := range tags {
		children = append(children, fuse.Dirent{
			Inode: t.fs.getInode(nTag, tag.ID),
			Name:  tag.Name,
		})
	}

	return children, nil
}

func (t *fsTags) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fusefs.Node, error) {
	id, err := t.fs.fsdb.AddTag(req.Name)
	if err != nil {
		return nil, err
	}

	return newFsTag(t.fs, id, req.Name), nil
}

func (t *fsTags) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	foundTag, err := t.fs.fsdb.GetTag(name)
	if err != nil {
		return nil, err
	}

	if foundTag == nil {
		return nil, fuse.ENOENT
	}

	return newFsTag(t.fs, foundTag.ID, foundTag.Name), nil
}

func (t *fsTags) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	if !req.Dir {
		return fuse.EIO
	}

	err := t.fs.fsdb.RemoveTag(req.Name)
	if err != nil {
		return err
	}

	return nil
}

type fsTag struct {
	node

	fs *DocFS

	id uint64
}

func newFsTag(fs *DocFS, id uint64, name string) *fsTag {
	t := &fsTag{
		fs: fs,
		id: id,
	}
	t.name = name
	return t
}

func (t *fsTag) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = t.fs.getInode(nTag, t.id)
	attr.Mode = os.ModeDir | 0755
	return nil
}
