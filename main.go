package main

//import (
//"bufio"
//"fmt"
//"log"
//"os"
//)

import (
	"bufio"
	"os"

	"github.com/PR2_MIA/analyzer"
	"github.com/PR2_MIA/commands"
	"github.com/PR2_MIA/execute"
	"github.com/PR2_MIA/systemCommands"
	"github.com/PR2_MIA/users"
	//"github.com/PR2_MIA/commands"
)

func main() {
	// LEO ENTRADA CON ESPACIOS HASTA ENCONTRAR UN SALTO DE LINEA
	in := bufio.NewReader(os.Stdin)
	inputCommand, _ := in.ReadString('\n')

	// OBTENGO EL ARBOL DE COMANDOS
	tree := analyzer.AnalyzerF(inputCommand, false)
	// BORRO LOS COMANDOS INNECESARIOS
	tree = tree[:2]

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
		} else if element.Name == "login" {
			login := users.LoginCmd{}
			login.AssignParameters(element)
			login.Login()
		} else if element.Name == "logout" {
			users.Logout()
		} else if element.Name == "mkgrp" {
			mkgrp := users.MkgrpCmd{}
			mkgrp.AssignParameters(element)
			mkgrp.Mkgrp()
		} else if element.Name == "rmgrp" {
			rmgrp := users.RmgrpCmd{}
			rmgrp.AssignParameters(element)
			rmgrp.Rmgrp()
		} else if element.Name == "mkuser" {
			mkuser := users.MkuserCmd{}
			mkuser.AssignParameters(element)
			mkuser.Mkuser()
		} else if element.Name == "rmusr" {
			rmuser := users.RmusrCmd{}
			rmuser.AssignParameters(element)
			rmuser.Rmusr()
		} else if element.Name == "pause" {
			commands.Pause()
		} else if element.Name == "exec" {
			exec := execute.ExecCmd{}
			exec.AssignParameters(element)
			exec.Exec()
		}
	}
}
