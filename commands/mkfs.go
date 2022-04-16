package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"

	"github.com/PR2_MIA/globals"
)

type MkfsCmd struct {
	Id   string
	Type string
}

func (cmd *MkfsCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "id" {
			cmd.Id = parameter.StringValue
		} else if parameter.Name == "type" {
			cmd.Type = parameter.StringValue
		}
	}
}

func (cmd *MkfsCmd) Mkfs() {
	if cmd.Id != "" {
		// BUSCO EL ID ENTRE LAS PARTICIONES MONTADAS
		mounted := globals.GlobalList.GetElement(cmd.Id)
		// SI NO LA ENCUENTRA LANZA UN ERROR
		if mounted.Id == "" {
			fmt.Println("Error: no existe particion montada con el id " + cmd.Id)
			return
		}

		// ABRO EL ARCHIVO
		file, err := os.OpenFile(mounted.Path, os.O_RDWR, 0777)

		if err != nil {
			log.Fatal("Error ", err)
			return
		}

		// INICIO SESION CON EL USUARIO ROOT EN EL GLOBAL USER
		globals.GlobalUser.Logged = 1
		globals.GlobalUser.Uid = "1"
		globals.GlobalUser.User_name = "root"
		globals.GlobalUser.Pwd = "123"
		globals.GlobalUser.Grp = "root"
		globals.GlobalUser.Id_partition = mounted.Id
		globals.GlobalUser.Gid = "1"

		// SI ES UN FORMATEO COMPLETO RELLAN DE CEROS EL ARCHIVO
		if cmd.Type == "full" || cmd.Type == "" {
			var temporal int8 = 0
			// RELLENO MI BUFFER CON CEROS
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, &temporal)

			for i := mounted.Start; i < (mounted.Start + mounted.Size); i++ {
				globals.WriteBytes(file, binario.Bytes())
			}
		}

		// CREO EL SUPERBLOQUE
		super_bloque := globals.SuperBloque{}
		inode := globals.InodeTable{}
		fileBlock := globals.FileBlock{}

		n := (mounted.Size - int(unsafe.Sizeof(super_bloque))) / (4 + int(unsafe.Sizeof(inode)) + (3 * (int(unsafe.Sizeof(fileBlock)))))
		fmt.Println(n)

		// INGRESO TODOS LOS VALORES DEL SUPERBLOQUE
		copy(super_bloque.Mnt_count[:], []byte(strconv.Itoa(1)))
		copy(super_bloque.Magic[:], []byte(strconv.Itoa(0xEF53)))
		copy(super_bloque.First_inode[:], []byte(strconv.Itoa(0)))
		copy(super_bloque.First_block[:], []byte(strconv.Itoa(0)))
		copy(super_bloque.Inodes_count[:], []byte(strconv.Itoa(n)))
		copy(super_bloque.Blocks_count[:], []byte(strconv.Itoa(3*n)))
		copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(n-2)))
		copy(super_bloque.Free_blocks_count[:], []byte(strconv.Itoa(((3 * n) - 2))))
		copy(super_bloque.Inode_size[:], []byte(strconv.Itoa(int(unsafe.Sizeof(inode)))))
		copy(super_bloque.Block_size[:], []byte(strconv.Itoa(int(unsafe.Sizeof(fileBlock)))))
		copy(super_bloque.Bm_inode_start[:], []byte(strconv.Itoa((mounted.Start + int(unsafe.Sizeof(super_bloque))))))
		copy(super_bloque.Filesystem_type[:], []byte(strconv.Itoa(2)))
		copy(super_bloque.Bm_block_start[:], []byte(strconv.Itoa((mounted.Start + int(unsafe.Sizeof(super_bloque)) + n))))
		copy(super_bloque.Inode_start[:], []byte(strconv.Itoa(((mounted.Start + int(unsafe.Sizeof(super_bloque)) + n) + (3 * n)))))
		copy(super_bloque.Block_start[:], []byte(strconv.Itoa((((mounted.Start + int(unsafe.Sizeof(super_bloque)) + n) + (3 * n)) + (n + int(unsafe.Sizeof(inode)))))))
		copy(super_bloque.Mnt_count[:], []byte(globals.GetDate()))

		var bitinodes = make([]byte, n)
		bitinodes[0] = '0'

	} else {
		fmt.Println("Error: el parametro id es obligatorio en el comando mkfs")
	}
}
