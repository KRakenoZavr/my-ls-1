package main

import (
	"log"
	"ml/flags"
	"ml/internal"
	"ml/utils"
	"os"
)

func main() {
	args := os.Args[1:]

	flagsArr, files := utils.GetArgs(args)

	flag, err := flags.NewFlags(flagsArr, len(files))
	if err != nil {
		log.Println(err, flagsArr)
		return
	}

	internal.Programm(files, flag)
}
