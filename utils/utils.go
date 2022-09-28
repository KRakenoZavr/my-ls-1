package utils

import (
	"io/fs"
	"os/user"
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
