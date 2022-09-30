package internal

import (
	"fmt"
	"io/fs"
	"ml/utils"
	"time"
)

type fileInfo struct {
	mode  fs.FileMode
	size  int64
	name  string
	isDir bool

	fullDate time.Time

	ownerGroup string
	ownerName  string

	isLink bool
	link   string

	blocks int64
}

func (f fileInfo) getMonthAsString() string {
	return f.fullDate.Month().String()[:3]
}

func (f fileInfo) formatMonth() string {
	return fmt.Sprintf("%s %2v %02d:%02d", f.getMonthAsString(), f.fullDate.Day(), f.fullDate.Hour(), f.fullDate.Minute())
}

func (f fileInfo) String() string {
	return f.GetName(false)
}

func (f fileInfo) SimplePrint() string {
	return f.name
}

// for print file name and link
func (f fileInfo) GetName(fullInfo bool) string {
	if f.isLink {
		if fullInfo {
			return fmt.Sprintf("%s%s%s -> %s%s%s", utils.Link, f.name, utils.Reset, utils.Dir, f.link, utils.Reset)
		}
		return fmt.Sprintf("%s%s%s", utils.Link, f.name, utils.Reset)
	}

	if f.isDir {
		return fmt.Sprintf("%s%s%s", utils.Dir, f.name, utils.Reset)
	}

	return f.name
}

// print all info of file
func (f fileInfo) FullPrint(fullInfo bool) string {
	return fmt.Sprintf("%s %s %s %4v %s %s", f.mode, f.ownerName, f.ownerGroup, f.size, f.formatMonth(), f.GetName(fullInfo))
}

type allFiles struct {
	files    []fileInfo
	fullInfo bool
}

func (f allFiles) getInfo() string {
	a := ""
	for _, l := range f.files {
		a += fmt.Sprintf("%s\n", l)
	}

	return a
}

func (f allFiles) getFullInfo() string {
	a := ""
	for _, l := range f.files {
		a += l.FullPrint(f.fullInfo) + "\n"
	}

	return a
}

// print func for allFiles struct
func (f allFiles) String() string {
	if !f.fullInfo {
		return f.getInfo()
	}

	return f.getFullInfo()
}
