package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/PR2_MIA/globals"
)

type MountCmd struct {
	Path string
	Name string
}

func (cmd *MountCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "name" {
			cmd.Name = parameter.StringValue
		}
	}
}

func (cmd *MountCmd) Mount() {
	if cmd.Name != "" {
		if cmd.Path != "" {
			// ABRO EL ARCHIVO
			file, err := os.OpenFile(cmd.Path, os.O_RDWR, 0777)

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
				return
			}

			primary_part := globals.Partition{}
			logic_part := globals.EBR{}
			isPrimary := false
			exist_partition := false

			for i := 0; i < 4; i++ {
				if globals.ByteToString(MBR.Partitions[i].Part_type[:]) == "e" {
					// SI ES UNA PARTICION EXTENDIDA BUSCO ENTRE TODAS SUS LOGICAS
					actualEbr := globals.EBR{}
					// LEO EL PRIMER EBR
					sizeEbr := int(unsafe.Sizeof(actualEbr))
					file.Seek(int64(globals.ByteToInt(MBR.Partitions[i].Part_start[:])), 0)
					dataEbr := globals.ReadBytes(file, sizeEbr)
					bufferebr := bytes.NewBuffer(dataEbr)
					errEbr := binary.Read(bufferebr, binary.BigEndian, &actualEbr)
					if errEbr != nil {
						log.Fatal("Error ", errEbr)
						return
					}
					if globals.ByteToString(MBR.Partitions[i].Part_name[:]) == cmd.Name {
						fmt.Println("Error: la particion " + cmd.Name + " no puede ser montada en ram debido a que es de tipo extendida")
						return
					}
					if globals.ByteToString(actualEbr.Part_name[:]) == cmd.Name {
						logic_part = actualEbr
						exist_partition = true
						/*fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
						return*/
					}
					// RECORRO TODOS LOS EBR
					for globals.ByteToInt(actualEbr.Part_next[:]) != -1 {
						size := int(unsafe.Sizeof(actualEbr))
						file.Seek(int64(globals.ByteToInt(actualEbr.Part_next[:])), 0)
						data := globals.ReadBytes(file, size)
						buffer := bytes.NewBuffer(data)
						err1 := binary.Read(buffer, binary.BigEndian, &actualEbr)
						if err1 != nil {
							log.Fatal("Error ", err1)
							return
						}
						if globals.ByteToString(actualEbr.Part_name[:]) == cmd.Name {
							logic_part = actualEbr
							exist_partition = true
							/*fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
							return*/
						}
					}
				} else {
					if globals.ByteToString(MBR.Partitions[i].Part_name[:]) == cmd.Name {
						primary_part = MBR.Partitions[i]
						isPrimary = true
						exist_partition = true
						/*fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
						return*/
					}
				}
			}

			if exist_partition {
				if isPrimary {
					disk_name := cmd.Path[strings.LastIndex(cmd.Path, "/")+1 : len(cmd.Path)]
					partition_name := globals.ByteToString(primary_part.Part_name[:])
					partition_start := globals.ByteToInt(primary_part.Part_start[:])
					type_partition := globals.ByteToString(primary_part.Part_type[:])
					fit_partition := globals.ByteToString(primary_part.Part_fit[:])
					size_partition := globals.ByteToInt(primary_part.Part_size[:])
					globals.GlobalList.MountPartition(disk_name, partition_name, partition_start, cmd.Path, type_partition, fit_partition, size_partition)
				} else {
					//fmt.Println(cmd.Name, globals.ByteToString(logic_part.Part_start[:]), globals.ByteToString(logic_part.Part_size[:]))
					disk_name := cmd.Path[strings.LastIndex(cmd.Path, "/")+1 : len(cmd.Path)]
					partition_name := globals.ByteToString(logic_part.Part_name[:])
					partition_start := globals.ByteToInt(logic_part.Part_start[:])
					//type_partition := globals.ByteToString(logic_part.Part_type[:])
					type_partition := "l"
					fit_partition := globals.ByteToString(logic_part.Part_fit[:])
					size_partition := globals.ByteToInt(logic_part.Part_size[:])
					globals.GlobalList.MountPartition(disk_name, partition_name, partition_start, cmd.Path, type_partition, fit_partition, size_partition)
				}
			} else {
				fmt.Println("Error: la particion " + cmd.Name + " no puede ser montada porque no existe en el disco " + cmd.Path)
			}
			file.Close()
		} else {
			fmt.Println("Error el parametro path es obligatorio en el comando mount")
		}
	} else {
		fmt.Println("Error el parametro name es obligatorio en el comando mount")
	}
}
