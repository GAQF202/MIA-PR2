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
	"github.com/PR2_MIA/systemCommands"
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
		copy(super_bloque.Bm_block_start[:], []byte(strconv.Itoa(globals.ByteToInt(super_bloque.Bm_inode_start[:])+n)))
		copy(super_bloque.Inode_start[:], []byte(strconv.Itoa(globals.ByteToInt(super_bloque.Bm_block_start[:])+(3*n))))
		copy(super_bloque.Block_start[:], []byte(strconv.Itoa(globals.ByteToInt(super_bloque.Inode_start[:])+(n*int(unsafe.Sizeof(inode))))))
		copy(super_bloque.Mnt_count[:], []byte(globals.GetDate()))

		// ESCRIBO EL SUPERBLOQUE EN EL DISCO
		file.Seek(int64(mounted.Start), 0)
		var bufferControl bytes.Buffer
		binary.Write(&bufferControl, binary.BigEndian, &super_bloque)
		globals.WriteBytes(file, bufferControl.Bytes())

		// CREACION DE BITMAPS
		// INODOS
		var bitinodes = make([]byte, n)
		for i := 0; i < n; i++ {
			bitinodes[i] = '0'
		}
		// OBTENGO LA POSICION DE BITMAP DE INODOS DEL SUPERBLOQUE
		bInodePos := globals.ByteToInt(super_bloque.Bm_inode_start[:])
		file.Seek(int64(bInodePos), 0)
		var bufferControlInodes bytes.Buffer
		binary.Write(&bufferControlInodes, binary.BigEndian, &bitinodes)
		globals.WriteBytes(file, bufferControlInodes.Bytes())

		// BLOQUES
		var bitblocks = make([]byte, (3 * n))
		for i := 0; i < 3*n; i++ {
			bitblocks[i] = '0'
		}
		// OBTENGO LA POSICION DE BITMAP DE INODOS DEL SUPERBLOQUE
		bBlockPos := globals.ByteToInt(super_bloque.Bm_block_start[:])
		file.Seek(int64(bBlockPos), 0)
		var bufferControlBlocks bytes.Buffer
		binary.Write(&bufferControlBlocks, binary.BigEndian, &bitblocks)
		globals.WriteBytes(file, bufferControlBlocks.Bytes())

		// CREO LA CARPETA RAIZ
		c := systemCommands.MkdirCmd{}
		c.Path = "/"
		c.P = "-p"
		c.Mkdir()

		d := systemCommands.MkfileCmd{}
		d.AnyText = "1,G,root\n1,U,root,root,123\n"
		d.Cont = ""
		d.Path = "/users.txt"
		d.R = "-r"
		d.Size = 0
		d.Mkfile()

		// CIERRO SESION
		globals.GlobalUser.Logged = -1
		globals.GlobalUser.Uid = ""
		globals.GlobalUser.User_name = ""
		globals.GlobalUser.Pwd = ""
		globals.GlobalUser.Grp = ""
		globals.GlobalUser.Id_partition = ""
		globals.GlobalUser.Gid = ""

	} else {
		fmt.Println("Error: el parametro id es obligatorio en el comando mkfs")
	}
}
