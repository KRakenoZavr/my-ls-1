package internal

import (
	"fmt"
	"io/fs"
	"log"
	"ml/flags"
	"ml/utils"
	"os"
	"sort"
	"time"
)

type fileDate struct {
	month string
	day   int
	time  string
}

type fileInfo struct {
	mode  fs.FileMode
	size  int64
	name  string
	isDir bool

	date     fileDate
	fullDate time.Time

	ownerGroup string
	ownerName  string

	isLink bool
	link   string
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
	return fmt.Sprintf("%s %v %s %s", f.mode, f.size, f.fullDate, f.GetName(fullInfo))
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

func Programm(files []string, flag *flags.Flag) {
	lots := false
	if len(files) > 1 {
		lots = true
	}
	// run programm for all arguments
	for i, l := range files {
		run(l, flag, lots, i == len(files)-1)
	}
}

func run(path string, flag *flags.Flag, lots, isLast bool) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	// check if dir
	fInfo, err := f.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	// if not dir, then run lsprog
	if !fInfo.IsDir() {
		lsprog([]fs.FileInfo{fInfo}, flag, lots, isLast)
		return
	}
	// if dir, then get all files in dir, and run lsprog
	files, err := f.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}

	// print which dir is it
	if lots {
		fmt.Printf("%s:\n", path)
	}

	lsprog(files, flag, lots, isLast)
}

func lsprog(files []fs.FileInfo, flag *flags.Flag, lots, isLast bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	fileInfos := allFiles{
		files:    []fileInfo{},
		fullInfo: true,
	}

	for _, v := range files {
		info := fileInfo{
			mode:     v.Mode(),
			size:     v.Size(),
			name:     v.Name(),
			isDir:    v.IsDir(),
			fullDate: v.ModTime(),
			isLink:   false,
		}

		// check of symlink
		if v.Mode()&os.ModeSymlink == os.ModeSymlink {
			// read link
			linkName, err := os.Readlink(v.Name())
			if err != nil {
				log.Println("error reading link:", err)
			} else {
				info.link = linkName
				info.isLink = true
			}
		}

		name1, name2, err := utils.GetOwnerFile(v)
		if err != nil {
			log.Println("error getting owner:", err)
		} else {
			info.ownerGroup = name1
			info.ownerName = name2
		}

		fileInfos.files = append(fileInfos.files, info)
	}

	// l - if not l, get only needed info
	if !flag.Contains("l") {
		fileInfos.fullInfo = false
	}

	// a - if not a, filter .
	if !flag.Contains("a") {
		fileInfos.files = RemoveDotFiles(fileInfos.files)
	}

	// t - sort by time
	if flag.Contains("t") {
		SortByTime(fileInfos.files)
	}

	// r - reverse
	if flag.Contains("r") {
		ReverseArray(fileInfos.files)
	}

	// print result
	fmt.Print(fileInfos)
	// print \n between multiple dirs
	if lots && !isLast {
		fmt.Println()
	}
	// R - recurse
}
