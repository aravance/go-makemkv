package mkv

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
)

type Device interface {
	Device() string
	Type() string
	Available() bool
}

type IsoDevice struct {
	path string
}

func (d *IsoDevice) Device() string {
	return d.path
}

func (d *IsoDevice) Type() string {
	return "iso"
}

func (d *IsoDevice) Available() bool {
	info, err := os.Stat(d.path)
	return err == nil && !info.IsDir()
}

type FileDevice struct {
	path string
}

func (d *FileDevice) Device() string {
	return d.path
}

func (d *FileDevice) Type() string {
	return "file"
}

func (d *FileDevice) Available() bool {
	info, err := os.Stat(d.path)
	return err == nil && info.IsDir()
}

type DevDevice struct {
	device string
}

func (d *DevDevice) Device() string {
	return "/dev/" + d.device
}

func (d *DevDevice) Type() string {
	return "disc"
}

func (d *DevDevice) Available() bool {
	_, err := os.Stat(d.Device())
	return err == nil || !errors.Is(err, fs.ErrNotExist)
}

type DiscDevice struct {
	id int
}

func (d *DiscDevice) Device() string {
	return strconv.Itoa(d.id)
}

func (d *DiscDevice) Type() string {
	return "dev"
}

func (d *DiscDevice) Available() bool {
	panic("not yet implemented")
}
