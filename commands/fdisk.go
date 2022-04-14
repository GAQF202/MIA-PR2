package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/PR2_MIA/globals"
)

type FdiskCmd struct {
	Size int
	Unit string
	Path string
	Type string
	Fit  string
	Name string
}

func (cmd *FdiskCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "size" {
			cmd.Size = parameter.IntValue
		} else if parameter.Name == "unit" {
			cmd.Unit = parameter.StringValue
		} else if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "type" {
			cmd.Type = parameter.StringValue
		} else if parameter.Name == "fit" {
			cmd.Fit = parameter.StringValue
		} else if parameter.Name == "name" {
			cmd.Name = parameter.StringValue
		}
	}
}

func (cmd *FdiskCmd) Fdisk() {
	if cmd.Size != -1 {
		if cmd.Name != "" {
			if cmd.Path != "" {
				// ABRO EL ARCHIVO
				file, err := os.Open(cmd.Path)

				if err != nil {
					log.Fatal("Error ", err)
					return
				}
				// LEO EL MBR
				MBR := globals.MBR{}
				size := int(unsafe.Sizeof(MBR))
				file.Seek(0, 0)
				data := globals.ReadBytes(file, size)
				buffer := bytes.NewBuffer(data)
				err1 := binary.Read(buffer, binary.BigEndian, &MBR)
				if err1 != nil {
					log.Fatal("Error ", err1)
				} else {
					fmt.Println(globals.ByteToString(MBR.Mbr_dsk_signature[:]))
				}

				// CALCULO DE MULTIPLICADOR PARA ASIGNAR ESPACIO A LA PARTICION
				multiplicator := 1024
				if cmd.Unit == "b" {
					multiplicator = 1
				} else if cmd.Unit == "m" {
					multiplicator = 1024 * 1024
				}
				fmt.Println(multiplicator)

				if cmd.Type == "" || cmd.Type == "p" || cmd.Type == "e" {
					fmt.Println("trabajar como primaria o extendida")
				} else {
					fmt.Println("trabajar com logica")
				}

				// CIERRO EL ARCHIVO
				file.Close()
			} else {
				fmt.Println("Error: el parametro path es obligatorio en fdisk")
			}
		} else {
			fmt.Println("Error: el parametro name es obligatorio en fdisk")
		}
	} else {
		fmt.Println("Error: el parametro size es obligatorio en fdisk")
	}
}
