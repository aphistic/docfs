package dfs

import (
	"fmt"
	"os"
	"strconv"

	"time"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type fsYear struct {
	node

	fs *DocFS

	Year uint64
}

func mergedDateKey(year uint64, month uint64, day uint64) uint64 {
	return (year * 10000) + (month * 100) + day
}

func newFsYear(fs *DocFS, year uint64) *fsYear {
	y := &fsYear{
		fs:   fs,
		Year: year,
	}
	y.inode = fs.getInode(nYear, uint64(year))
	y.name = fmt.Sprintf("%d", year)
	return y
}

func (y *fsYear) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = y.inode
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (y *fsYear) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var children []fuse.Dirent

	months, err := y.fs.fsdb.GetMonths(y.Year)
	if err != nil {
		return nil, err
	}

	for _, month := range months {
		children = append(children, fuse.Dirent{
			Name:  fmt.Sprintf("%02d", month),
			Inode: y.fs.getInode(nMonth, mergedDateKey(y.Year, month, 0)),
		})
	}

	return children, nil
}

func (y *fsYear) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	month, err := strconv.ParseUint(name, 10, 0)
	if err != nil {
		return nil, fuse.ENOENT
	}

	month, err = y.fs.fsdb.GetMonth(y.Year, month)
	if err != nil {
		return nil, fuse.ENOENT
	}

	return newFsMonth(y.fs, y.Year, month), nil
}

func (y *fsYear) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fusefs.Node, error) {
	month, err := strconv.ParseUint(req.Name, 10, 0)
	if err != nil {
		return nil, fuse.EPERM
	}

	if month < 1 || month > 12 {
		return nil, fuse.EPERM
	}

	err = y.fs.fsdb.AddMonth(y.Year, month)
	if err != nil {
		return nil, err
	}

	return newFsMonth(y.fs, y.Year, month), nil
}

func (y *fsYear) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	if !req.Dir {
		return fuse.EIO
	}

	month, err := strconv.ParseUint(req.Name, 10, 0)
	if err != nil {
		return fuse.EIO
	}

	err = y.fs.fsdb.RemoveMonth(y.Year, month)
	if err != nil {
		return err
	}

	return nil
}

type fsMonth struct {
	node

	fs *DocFS

	Year  uint64
	Month uint64
}

func newFsMonth(fs *DocFS, year uint64, month uint64) *fsMonth {
	m := &fsMonth{
		fs:    fs,
		Year:  year,
		Month: month,
	}
	m.inode = fs.getInode(nMonth, mergedDateKey(m.Year, m.Month, 0))
	m.name = fmt.Sprintf("%02d", month)
	return m
}

func (m *fsMonth) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = m.inode
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (m *fsMonth) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var children []fuse.Dirent

	days, err := m.fs.fsdb.GetDays(m.Year, m.Month)
	if err != nil {
		return nil, err
	}

	for _, day := range days {
		children = append(children, fuse.Dirent{
			Name:  fmt.Sprintf("%02d", day),
			Inode: m.fs.getInode(nDay, mergedDateKey(m.Year, m.Month, day)),
		})
	}

	return children, nil
}

func (m *fsMonth) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	day, err := strconv.ParseUint(name, 10, 0)
	if err != nil {
		return nil, fuse.ENOENT
	}

	day, err = m.fs.fsdb.GetDay(m.Year, m.Month, day)
	if err != nil {
		return nil, fuse.ENOENT
	}

	return newFsDay(m.fs, m.Year, m.Month, day), nil
}

func (m *fsMonth) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fusefs.Node, error) {
	day, err := strconv.ParseUint(req.Name, 10, 0)
	if err != nil {
		return nil, fuse.EPERM
	}

	nextMonth := time.Date(int(m.Year), time.Month(m.Month+1), 1, 0, 0, 0, 0, time.UTC)
	lastDay := nextMonth.Add(-1 * time.Hour)

	if day < 1 || int(day) > lastDay.Day() {
		return nil, fuse.EPERM
	}

	err = m.fs.fsdb.AddDay(m.Year, m.Month, day)
	if err != nil {
		return nil, err
	}

	return newFsDay(m.fs, m.Year, m.Month, day), nil
}

func (m *fsMonth) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	if !req.Dir {
		return fuse.EIO
	}

	day, err := strconv.ParseUint(req.Name, 10, 0)
	if err != nil {
		return fuse.EIO
	}

	err = m.fs.fsdb.RemoveDay(m.Year, m.Month, day)
	if err != nil {
		return err
	}

	return nil
}

type fsDay struct {
	node

	fs *DocFS

	Year  uint64
	Month uint64
	Day   uint64
}

func newFsDay(fs *DocFS, year uint64, month uint64, day uint64) *fsDay {
	d := &fsDay{
		fs:    fs,
		Year:  year,
		Month: month,
		Day:   day,
	}
	d.inode = fs.getInode(nDay, mergedDateKey(d.Year, d.Month, day))
	d.name = fmt.Sprintf("%02d", day)
	return d
}

func (d *fsDay) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = d.inode
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (d *fsDay) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var children []fuse.Dirent
	return children, nil
}

func (d *fsDay) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	return nil, fuse.ENOENT

}

func (d *fsDay) Create(ctx context.Context, req *fuse.CreateRequest, res *fuse.CreateResponse) (fusefs.Node, fusefs.Handle, error) {
	doc := newFsDoc(d.fs, 0, req.Name)

	return doc, nil, nil
}
