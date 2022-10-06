package internal

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"syscall"
	"time"

	"ml/utils"
)

type fileInfo struct {
	mode     fs.FileMode
	size     int64
	name     string
	isDir    bool
	fullPath string

	fullDate time.Time

	ownerGroup string
	ownerName  string

	isLink bool
	link   string

	blocks int64

	hardLinks int
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
func (f fileInfo) FullPrint(fullInfo bool, maxHL, maxSize, maxOG, maxOU int) string {
	return fmt.Sprintf("%s %*v %-*s %-*s %*v %s %s", f.mode, maxHL, f.hardLinks, maxOU, f.ownerName, maxOG, f.ownerGroup, maxSize, f.size, f.formatMonth(), f.GetName(fullInfo))
}

// add block info
func (f *fileInfo) AddBlocks() {
	var fileStats syscall.Stat_t
	err := syscall.Lstat(f.fullPath, &fileStats)
	if err != nil {
		log.Println(err)
	} else {
		f.blocks = fileStats.Blocks
	}
}

// add symlink info
func (f *fileInfo) Symlink() {
	if f.mode&os.ModeSymlink == os.ModeSymlink {
		// read link
		linkName, err := os.Readlink(f.fullPath)
		if err != nil {
			log.Println("error reading link:", err)
		} else {
			f.link = linkName
			f.isLink = true
		}
	}
}

// add symlink info
func (f *fileInfo) FileOwner(file fs.FileInfo) {
	name1, name2, err := utils.GetOwnerFile(file)
	if err != nil {
		log.Println("error getting owner:", err)
	} else {
		f.ownerGroup = name1
		f.ownerName = name2
	}
}

type allFiles struct {
	files    []*fileInfo
	fullInfo bool
	isDir    bool
	path     string
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
	f.hardLink()
	if f.isDir {
		a = fmt.Sprintf("total %v\n", f.totalBlocks())
	}

	maxHL := 0
	maxSize := 0
	maxOG := 0
	maxOU := 0
	for _, l := range f.files {
		lenHL := len(fmt.Sprintf("%v", l.hardLinks))
		if maxHL < lenHL {
			maxHL = lenHL
		}

		lenSize := len(fmt.Sprintf("%v", l.size))
		if maxSize < lenSize {
			maxSize = lenSize
		}

		lenOG := len(l.ownerGroup)
		if maxOG < lenOG {
			maxOG = lenOG
		}

		lenOU := len(l.ownerName)
		if maxOU < lenOU {
			maxOU = lenOU
		}
	}

	for _, l := range f.files {
		a += l.FullPrint(f.fullInfo, maxHL, maxSize, maxOG, maxOU) + "\n"
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

// get total blocks
func (f allFiles) totalBlocks() (c int64) {
	for _, l := range f.files {
		c += l.blocks
	}

	return int64(math.Round(float64(c) * physicalBlockSize / lsBlockSize))
}

func (f allFiles) hardLink() {
	for _, l := range f.files {
		if !l.isDir {
			l.hardLinks = 1
			continue
		}

		f, err := os.Open(l.fullPath)
		if err != nil {
			log.Println(err)
			l.hardLinks = 1
			continue
		}
		defer f.Close()

		files, err := f.Readdir(0)
		if err != nil {
			log.Println(err)
			l.hardLinks = 1
			continue
		}

		c := 2
		for _, l2 := range files {
			if l2.IsDir() {
				c++
			}
		}
		l.hardLinks = c
	}
}

func (f *allFiles) curDir() {
	v, err := os.Lstat(f.path)
	if err != nil {
		log.Println(err)
		return
	}

	info := &fileInfo{
		mode:     v.Mode(),
		size:     v.Size(),
		name:     ".",
		isDir:    v.IsDir(),
		fullDate: v.ModTime(),
		isLink:   false,
		fullPath: utils.GetPathToLink(f.path, v.Name()),
		blocks:   0,
	}

	info.AddBlocks()
	info.Symlink()
	info.FileOwner(v)

	f.files = append([]*fileInfo{info}, f.files...)
}

func (f *allFiles) parentDir() {
	v, err := os.Lstat(f.path + "/..")
	if err != nil {
		log.Println(err)
		return
	}

	info := &fileInfo{
		mode:     v.Mode(),
		size:     v.Size(),
		name:     v.Name(),
		isDir:    v.IsDir(),
		fullDate: v.ModTime(),
		isLink:   false,
		fullPath: utils.GetPathToLink(f.path+"/..", v.Name()),
		blocks:   0,
	}

	info.AddBlocks()
	info.Symlink()
	info.FileOwner(v)

	f.files = append([]*fileInfo{info}, f.files...)
}

func (f *allFiles) dotDirs() {
	f.parentDir()
	f.curDir()
}
