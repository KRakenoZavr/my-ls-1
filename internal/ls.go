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
	return f.name
}

func (f fileInfo) SimplePrint() string {
	return f.name
}

func (f fileInfo) GetName() string {
	if f.isLink {
		return fmt.Sprintf("%s %s %s -> %s %s %s", utils.Link, f.name, utils.Reset, utils.Dir, f.link, utils.Reset)
	}

	if f.isDir {
		return fmt.Sprintf("%s %s %s", utils.Dir, f.name, utils.Reset)
	}

	return f.name
}

func (f fileInfo) FullPrint() string {
	return fmt.Sprintf("%s %v %s %s", f.mode, f.size, f.fullDate, f.GetName())
}

type allFiles struct {
	files    []fileInfo
	fullInfo bool
}

func (f allFiles) String() string {
	if !f.fullInfo {
		return fmt.Sprintf("%s", f.files)
	}
	a := ""
	for _, l := range f.files {
		a += l.FullPrint() + "\n"
	}
	return a
}

func Programm(files []string, flag *flags.Flag) {
	for _, l := range files {
		run(l, flag)
	}
}

func run(path string, flag *flags.Flag) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	fInfo, err := f.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	if !fInfo.IsDir() {
		lsprog([]fs.FileInfo{fInfo}, flag)
		return
	}

	files, err := f.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}

	lsprog(files, flag)
}

func lsprog(files []fs.FileInfo, flag *flags.Flag) {
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

		// fmt.Println(v.Mode(), v.Size(), v.ModTime(), v.Name(), v.IsDir())
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

	fmt.Println(fileInfos)
	// R - recurse
}
