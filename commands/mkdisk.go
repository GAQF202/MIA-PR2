package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PR2_MIA/globals"
)

type MkdiskCmd struct {
	Size int
	Fit  string
	Unit string
	Path string
}

func (cmd *MkdiskCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "size" {
			cmd.Size = parameter.IntValue
		} else if parameter.Name == "fit" {
			cmd.Fit = parameter.StringValue
		} else if parameter.Name == "unit" {
			cmd.Unit = parameter.StringValue
		} else if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		}
	}
}

func (cmd *MkdiskCmd) Mkdisk() {
	if cmd.Size != -1 {
		// CREA LAS CARPETAS PADRE
		parent_path := cmd.Path[:strings.LastIndex(cmd.Path, "/")]
		if err := os.MkdirAll(parent_path, 0700); err != nil {
			log.Fatal(err)
		}

		// CREA EL ARCHIVO
		disk_file, err := os.Create(cmd.Path)
		if err != nil {
			log.Fatal(err)
		}

		// RELLENA EL ARCHIVO CON CEROS
		multiplicator := 1024
		if cmd.Unit == "k" {
			multiplicator = 1
		}
		var temporal int8 = 0
		// RELLENO MI BUFFER CON CEROS
		var binario bytes.Buffer
		for i := 0; i < 1024; i++ {
			binary.Write(&binario, binary.BigEndian, &temporal)
		}

		for i := 0; i < cmd.Size*multiplicator; i++ {
			globals.WriteBytes(disk_file, binario.Bytes())
		}

		// CREO EL MBR
		MBR := globals.MBR{}

		// INICIALIZO LAS PARTICIONES DEL MBR
		for i := 0; i < 4; i++ {
			copy(MBR.Partitions[i].Part_name[:], "")
			copy(MBR.Partitions[i].Part_status[:], "0")
			copy(MBR.Partitions[i].Part_type[:], "P")
			copy(MBR.Partitions[i].Part_start[:], "-1")
			copy(MBR.Partitions[i].Part_size[:], "-1")
			copy(MBR.Partitions[i].Part_fit[:], []byte(cmd.Fit))
		}

		// ASIGNACION DE ATRIBUTOS DEL MBR
		copy(MBR.Dsk_fit[:], []byte(cmd.Fit))
		copy(MBR.Mbr_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator*1024)))
		copy(MBR.Mbr_dsk_signature[:], []byte(strconv.Itoa(globals.GetRandom())))
		copy(MBR.Mbr_fecha_creacion[:], []byte(globals.GetDate()))

		// ESCRIBO EL MBR EN EL DISCO
		disk_file.Seek(0, 0)
		var bufferControl bytes.Buffer
		binary.Write(&bufferControl, binary.BigEndian, &MBR)
		globals.WriteBytes(disk_file, bufferControl.Bytes())

		/*disk_file.Seek(0, 0)
		temp := globals.MBR{}
		var size int = int(unsafe.Sizeof(temp))

		data := globals.ReadBytes(disk_file, size)
		buf := bytes.NewBuffer(data)
		e := binary.Read(buf, binary.BigEndian, &temp)

		if e != nil {
			log.Fatal("binary.Read has failed", e)
		} else {
			fmt.Println(string(temp.Mbr_size[:]))
		}*/
		// CIERRO EL ARCHIVO
		disk_file.Close()

	} else {
		fmt.Println("Error: el parametro size es obligatorio en mkdisk")
	}
}
