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

const (
	physicalBlockSize = 512
	lsBlockSize       = 1024
)

func Programm(files []string, flag *flags.Flag) {
	lots := false
	if len(files) > 1 {
		lots = true
	}

	if flag.Contains("R") {
		lots = true
	}

	sortedFiles := utils.SortFiles(files)

	// run programm for all arguments
	for i, l := range sortedFiles {
		run(l, flag, lots, i == 0)
	}
}

func run(path string, flag *flags.Flag, lots, isFirst bool) {
	// check if dir
	fInfo, err := os.Lstat(path)
	if err != nil {
		log.Println(err)
		return
	}

	// if not dir, then run lsprog
	if !fInfo.IsDir() {
		lsprog([]fs.FileInfo{fInfo}, flag, lots, false, path)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	// if dir, then get all files in dir, and run lsprog
	files, err := f.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}

	// print \n between multiple dirs
	if lots && !isFirst {
		fmt.Println()
	}
	// print which dir is it
	if lots {
		fmt.Printf("%s:\n", path)
	}

	lsprog(files, flag, lots, true, path)
}

func lsprog(files []fs.FileInfo, flag *flags.Flag, lots, isDir bool, path string) {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	fileInfos := allFiles{
		files:    []*fileInfo{},
		fullInfo: true,
		isDir:    isDir,
		path:     path,
	}

	for _, v := range files {
		info := &fileInfo{
			mode:     v.Mode(),
			size:     v.Size(),
			name:     v.Name(),
			isDir:    v.IsDir(),
			fullDate: v.ModTime(),
			isLink:   false,
			fullPath: utils.GetPathToLink(path, v.Name()),
			blocks:   0,
		}

		// blocks
		if flag.Contains("l") {
			info.AddBlocks()
		}

		// check of symlink
		info.Symlink()
		info.FileOwner(v)

		fileInfos.files = append(fileInfos.files, info)
	}

	// l - if not l, get only needed info
	if !flag.Contains("l") {
		fileInfos.fullInfo = false
	}

	// a - if not a, filter .
	if !flag.Contains("a") {
		fileInfos.files = RemoveDotFiles(fileInfos.files)
	} else {
		fileInfos.dotDirs()
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

	// R - recurse
	if flag.Contains("R") {
		for _, l := range fileInfos.files {
			if l.isDir {
				run(fmt.Sprintf("%s/%s", path, l.name), flag, true, false)
			}
		}
	}
}
