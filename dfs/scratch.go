package dfs

import (
	"fmt"

	"os"

	"hash"

	"crypto/sha256"

	"encoding/hex"

	"bufio"

	"bazil.org/fuse"
	"golang.org/x/net/context"
)

type scratchDoc struct {
	fs *DocFS

	id    uint64
	inode uint64

	hasher hash.Hash

	file    *os.File
	fileBuf *bufio.Writer
}

func newScratchDoc(fs *DocFS, id uint64) *scratchDoc {
	return &scratchDoc{
		fs: fs,

		id:    id,
		inode: fs.getInode(nScratch, id),

		hasher: sha256.New(),
	}
}

func (s *scratchDoc) Open() error {
	file, err := s.fs.openScratch(s.id)
	if err != nil {
		return err
	}

	s.file = file
	s.fileBuf = bufio.NewWriter(s.file)

	return nil
}

func (s *scratchDoc) Close() error {
	hash := s.hasher.Sum(nil)
	hashStr := hex.EncodeToString(hash)
	fmt.Printf("Closing scratch %d with hash %s\n", s.id, hashStr)

	return s.file.Close()
}

func (s *scratchDoc) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = s.inode
	attr.Mode = 0755
	return nil
}

func (s *scratchDoc) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	writeN, err := s.fileBuf.Write(req.Data)
	if err != nil {
		return err
	}

	hashN, err := s.hasher.Write(req.Data)
	if err != nil {
		return err
	}

	if writeN != hashN {
		return ErrHashFailure
	}

	resp.Size = writeN
	return nil
}

func (s *scratchDoc) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	return s.fileBuf.Flush()
}

func (s *scratchDoc) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	return s.Close()
}
