package utils

import (
	"io/fs"
	"log"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// TODO check files
// TODO check flags?
func GetArgs(args []string) (flags []string, files []string) {
	for _, l := range args {
		if strings.HasPrefix(l, "-") {
			flags = append(flags, l)
		} else {
			files = append(files, l)
		}
	}

	// if no args, then ls current dir
	if len(files) == 0 {
		files = append(files, ".")
	}

	return
}

func GetOwnerFile(file fs.FileInfo) (OwnerGroup string, OwnerName string, err error) {
	stat := file.Sys().(*syscall.Stat_t)

	uid := stat.Uid
	gid := stat.Gid

	u := strconv.FormatUint(uint64(uid), 10)
	g := strconv.FormatUint(uint64(gid), 10)

	usr, err := user.LookupId(u)
	if err != nil {
		return
	}
	group, err := user.LookupGroupId(g)
	if err != nil {
		return
	}

	return group.Name, usr.Username, err
}

func SortFiles(args []string) []string {
	dirs := make([]string, 0)
	files := make([]string, 0)

	for _, l := range args {
		f, err := os.Open(l)
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()
		// check if dir
		fInfo, err := f.Stat()
		if err != nil {
			log.Println(err)
			continue
		}

		// if not dir, then run lsprog
		if !fInfo.IsDir() {
			files = append(files, l)
		} else {
			dirs = append(dirs, l)
		}
	}

	sort.SliceStable(files, func(i, j int) bool {
		return strings.ToLower(files[i]) < strings.ToLower(files[j])
	})

	sort.SliceStable(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i]) < strings.ToLower(dirs[j])
	})

	return append(files, dirs...)
}

func GetPathToLink(path, filename string) string {
	var linkPath string

	if strings.HasSuffix(path, "/") {
		linkPath = path + filename
	} else if strings.HasSuffix(path, filename) {
		linkPath = path
	} else {
		linkPath = path + "/" + filename
	}

	return linkPath
}
