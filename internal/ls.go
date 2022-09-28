package internal

import (
	"fmt"
	"io/fs"
	"log"
	"ml/flags"
	"ml/utils"
	"os"
	"sort"
)

func Programm(files []string, flag *flags.Flag) {
	lots := false
	if len(files) > 1 {
		lots = true
	}

	sortedFiles := utils.SortFiles(files)

	// run programm for all arguments
	for i, l := range sortedFiles {
		run(l, flag, lots, i == len(files)-1)
	}
}

func run(path string, flag *flags.Flag, lots, isLast bool) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
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
