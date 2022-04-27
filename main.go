package main

//import (
//"bufio"
//"fmt"
//"log"
//"os"
//)

import (
	"github.com/PR2_MIA/analyzer"
	"github.com/PR2_MIA/commands"
	"github.com/PR2_MIA/systemCommands"
	//"github.com/PR2_MIA/commands"
)

func main() {
	// OBTENGO EL ARBOL DE COMANDOS
	tree := analyzer.AnalyzerF()

	for _, element := range tree {
		if element.Name == "mkdisk" {
			mkdisk := commands.MkdiskCmd{}
			mkdisk.AssignParameters(element)
			mkdisk.Mkdisk()
		} else if element.Name == "rmdisk" {
			rmdisk := commands.RmdiskCmd{}
			rmdisk.AssignParameters(element)
			rmdisk.Rmdisk()
		} else if element.Name == "fdisk" {
			fdisk := commands.FdiskCmd{}
			fdisk.AssignParameters(element)
			fdisk.Fdisk()
		} else if element.Name == "mount" {
			mount := commands.MountCmd{}
			mount.AssignParameters(element)
			mount.Mount()
		} else if element.Name == "mkfs" {
			mkfs := commands.MkfsCmd{}
			mkfs.AssignParameters(element)
			mkfs.Mkfs()
		} else if element.Name == "mkdir" {
			mkdir := systemCommands.MkdirCmd{}
			mkdir.AssignParameters(element)
			mkdir.Mkdir()
		} else if element.Name == "mkfile" {
			mkfile := systemCommands.MkfileCmd{}
			mkfile.AssignParameters(element)
			mkfile.Mkfile()
		} else if element.Name == "rep" {
			rep := commands.RepCmd{}
			rep.AssignParameters(element)
			rep.Rep()
		} else if element.Name == "comment" {
			comment := commands.Comment{}
			comment.AssignParameters(element)
			comment.ShowComment()
		}
	}
}
