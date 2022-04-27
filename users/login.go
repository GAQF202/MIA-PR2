package users

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/PR2_MIA/globals"
	"github.com/PR2_MIA/read"
)

type LoginCmd struct {
	Usuario  string
	Password string
	Id       string
}

func (cmd *LoginCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "usuario" {
			cmd.Usuario = parameter.StringValue
		} else if parameter.Name == "password" {
			cmd.Password = parameter.StringValue
		} else if parameter.Name == "id" {
			cmd.Id = parameter.StringValue
		}
	}
}

func (cmd *LoginCmd) Login() {

	if cmd.Id != "" {
		if cmd.Password != "" {
			if cmd.Usuario != "" {
				// VARIABLE CON TODA LA INFORMACION DE LA PARTICION MONTADA
				partition_m := globals.GlobalList.GetElement(cmd.Id)

				// ABRO EL ARCHIVO
				file, err := os.OpenFile(partition_m.Path, os.O_RDWR, 0777)
				// VERIFICACION DE ERROR AL ABRIR EL ARCHIVO
				if err != nil {
					log.Fatal("Error ", err)
					return
				}

				// LEO EL SUPERBLOQUE
				super_bloque := globals.SuperBloque{}
				super_bloque = read.ReadSuperBlock(file, partition_m.Start)

				// LEO EL PRIMER INODO QUE ES EL QUE CONTIENE EL ARCHIVO DE USUARIOS
				users_inode := globals.InodeTable{}
				users_inode = read.ReadInode(file, globals.ByteToInt(super_bloque.Inode_start[:])+int(unsafe.Sizeof(users_inode)))

				archive_block := globals.ArchiveBlock{}

				users_archive_content := ""

				for block_i := 0; block_i < 16; block_i++ {
					if users_inode.Block[block_i] != '1' {
						archive_block = read.ReadArchiveBlock(file, globals.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[block_i])*int(unsafe.Sizeof(archive_block))))
						users_archive_content += globals.ByteToString(archive_block.Content[:])
					}
				}

				// ALMACENO TODOS LOS GRUPOS Y USUARIOS SEPARADOS POR UN SALTO
				all := strings.Split(users_archive_content, "\n")
				// ARREGLOS PARA GUARDAR LOS GRUPOS Y USUARIOS POR SEPARADO
				var groups = make([]globals.Group, 0)
				var users = make([]globals.User, 0)

				// RECORRO TODOS LOS USUAIROS Y GRUPOSY LOS SEPARO
				for i := 0; i < len(all); i++ {
					if all[i] != "" {
						temp := strings.Split(all[i], ",")
						if temp[1] == "G" {
							groups = append(groups, globals.Group{temp[0], temp[1], temp[2]})
						} else if temp[1] == "U" {
							users = append(users, globals.User{temp[0], temp[1], temp[2], temp[3], temp[4]})
						}
					}
				}

				exist_user_in := false
				for i := 0; i < len(users); i++ {
					if cmd.Usuario == users[i].User && cmd.Password == users[i].Password {
						exist_user_in = true
						// INICIO SESION
						globals.GlobalUser.Logged = 1
						globals.GlobalUser.Uid = users[i].Uid
						globals.GlobalUser.User_name = users[i].User
						globals.GlobalUser.Pwd = users[i].Password
						globals.GlobalUser.Grp = users[i].Group
						globals.GlobalUser.Id_partition = partition_m.Id
						globals.GlobalUser.Gid = users[i].Group
						break
					}
				}

				if !exist_user_in {
					fmt.Println("Error: el usuario o la contraseña en login no son correctas")
				}
			} else {
				fmt.Println("Error: el parametro usuario es obligatorio en el comando login")
			}
		} else {
			fmt.Println("Error: el parametro password es obligatorio en el comando login")
		}
	} else {
		fmt.Println("Error: el parametro id es obligatorio en el comando login")
	}

}
