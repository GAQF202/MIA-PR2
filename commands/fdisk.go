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
				// ARRAY PARA GUARDAR LOS ESPACIOS LIBRES
				var avalaibleSpaces = make([]globals.VoidSpace, 0)

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
						if globals.ByteToString(actualEbr.Part_name[:]) == cmd.Name || globals.ByteToString(MBR.Partitions[i].Part_name[:]) == cmd.Name {
							fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
							return
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
								fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
								return
							}
						}
					} else {
						if globals.ByteToString(MBR.Partitions[i].Part_name[:]) == cmd.Name {
							fmt.Println("Error: la particion con el nombre " + cmd.Name + " ya existe en el disco " + cmd.Path)
							return
						}
					}
				}

				// CALCULO DE MULTIPLICADOR PARA ASIGNAR ESPACIO A LA PARTICION
				multiplicator := 1024
				if cmd.Unit == "b" {
					multiplicator = 1
				} else if cmd.Unit == "m" {
					multiplicator = 1024 * 1024
				}
				//fmt.Println(multiplicator)

				// ORDENA TODAS LAS PARTICIONES DE MENOR A MAYOR
				globals.BubbleSort(MBR.Partitions[:])
				// CALCULO DE TIPO DE PARTICIONES
				totalPartitions := 0
				totalExtended := 0
				extendedPartition := globals.Partition{}

				for i := 0; i < 4; i++ {
					//fmt.Println(globals.ByteToString(MBR.Partitions[i].Part_size[:]), globals.ByteToString(MBR.Partitions[i].Part_start[:]))
					// CUENTA CUANTAS PARTICIONES HAY DENTRO DEL MBR
					if globals.ByteToString(MBR.Partitions[i].Part_size[:]) != "-1" {
						totalPartitions++
					}

					if globals.ByteToString(MBR.Partitions[i].Part_type[:]) == "e" {
						extendedPartition = MBR.Partitions[i]
						totalExtended++
					}
				}

				if cmd.Type == "" || cmd.Type == "p" || cmd.Type == "e" {

					/*for i := 0; i < 4; i++ {
						fmt.Println(globals.ByteToString(MBR.Partitions[i].Part_start[:]))
					}*/
					// CALCULO LOS ESPACIOS VACIOS ENTRE PARTICIONES PRIMARIAS Y EXTENDIDAS
					if totalPartitions == 4 {
						fmt.Println("Error: la particion " + cmd.Name + " no se pude montar porque la suma de particiones extendidas y primarias llego a su límite")
						return
					} else if totalExtended > 0 && cmd.Type == "e" {
						fmt.Println("Error: la particion " + cmd.Name + " no se pude montar porque solamente puede existir una particion extendida en el disco")
						return
					}

					if totalPartitions != 0 {
						for i := 0; i < 4; i++ {
							if i == 0 {
								tmpSpace := globals.VoidSpace{}
								tmpSpace.Size = globals.ByteToInt(MBR.Partitions[i].Part_start[:]) - int(unsafe.Sizeof(MBR)) - 2
								tmpSpace.Start = int(unsafe.Sizeof(MBR)) + 1
								avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
							} else if i == 3 {
								tmpSpace := globals.VoidSpace{}
								tmpSpace.Size = globals.ByteToInt(MBR.Mbr_size[:]) - (globals.ByteToInt(MBR.Partitions[i].Part_size[:]) + globals.ByteToInt(MBR.Partitions[i].Part_start[:])) - 1
								tmpSpace.Start = globals.ByteToInt(MBR.Partitions[i].Part_size[:]) + globals.ByteToInt(MBR.Partitions[i].Part_start[:]) + 1
								avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
							} else {
								tmpSpace := globals.VoidSpace{}
								tmpSpace.Size = globals.ByteToInt(MBR.Partitions[i].Part_start[:]) - (globals.ByteToInt(MBR.Partitions[i-1].Part_size[:]) + globals.ByteToInt(MBR.Partitions[i-1].Part_start[:])) - 2
								tmpSpace.Start = globals.ByteToInt(MBR.Partitions[i-1].Part_size[:]) + globals.ByteToInt(MBR.Partitions[i-1].Part_start[:]) + 1
								avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
							}
						}
					} else {
						tmpSpace := globals.VoidSpace{}
						tmpSpace.Size = globals.ByteToInt(MBR.Mbr_size[:]) - (int(unsafe.Sizeof(MBR)) + 1)
						tmpSpace.Start = int(unsafe.Sizeof(MBR)) + 1
						avalaibleSpaces = append(avalaibleSpaces, tmpSpace)
					}
					/*fmt.Println("-------------------------------")
					for i := 0; i < len(avalaibleSpaces); i++ {
						fmt.Println(avalaibleSpaces[i])
					}
					fmt.Println("-------------------------------")*/
					// ORDENO LOS ESPACIOS VACIOS DE MENOR A MAYOR TAMANIO
					globals.SortFreeSpaces(avalaibleSpaces[:])

					// VARIABLE PARA GUARDAR DONDE INICIA LA PARTICION CREADA
					selectVoidSpace := -1
					for i := 0; i < len(avalaibleSpaces); i++ {
						//fmt.Println("Estoo", avalaibleSpaces[i].Size, cmd.Size*multiplicator, avalaibleSpaces[i].Start)
						if avalaibleSpaces[i].Size >= cmd.Size*multiplicator {
							selectVoidSpace = avalaibleSpaces[i].Start
							break
						}
					}
					if selectVoidSpace != -1 {
						for i := 0; i < 4; i++ {
							if globals.ByteToInt(MBR.Partitions[i].Part_size[:]) == -1 {
								fit := ""
								ptype := ""
								if cmd.Fit == "" {
									fit = "wf"
								}
								if cmd.Type == "" {
									ptype = "p"
								} else {
									ptype = cmd.Type
								}

								// SI ES UNA PARTICION EXTENDIDA CREO LA CABECERA
								if ptype == "e" {
									// ESCRIBO EL EBR INICIAL
									EBR := globals.EBR{}
									copy(EBR.Part_status[:], "0")
									copy(EBR.Part_fit[:], "wf")
									copy(EBR.Part_start[:], "-1")
									copy(EBR.Part_size[:], "-1")
									copy(EBR.Part_next[:], "-1")
									copy(EBR.Part_name[:], "")
									// ME POSICIONO AL INICIO DE LA PARTICION EXTENDIDA
									file.Seek(int64(selectVoidSpace), 0)
									// ESCRIBO EL PRIMER EBR EN EL DISCO
									var bufferControl bytes.Buffer
									binary.Write(&bufferControl, binary.BigEndian, &EBR)
									globals.WriteBytes(file, bufferControl.Bytes())
								}

								copy(MBR.Partitions[i].Part_name[:], []byte(cmd.Name))
								copy(MBR.Partitions[i].Part_fit[:], []byte(fit))
								copy(MBR.Partitions[i].Part_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator)))
								copy(MBR.Partitions[i].Part_start[:], []byte(strconv.Itoa(selectVoidSpace)))
								copy(MBR.Partitions[i].Part_type[:], []byte(ptype))
								// IMPRIMO EN CONSOLA LOS DATOS DE LA PARTICION CREADA
								globals.PrintPartition(MBR.Partitions[i])
								break
							}
						}
					} else {
						fmt.Println("Error: la particion " + cmd.Name + " no cabe en el disco " + cmd.Path)
					}
				} else {
					if totalExtended == 1 {
						fit := ""
						if cmd.Fit == "" {
							fit = "wf"
						} else {
							fit = cmd.Fit
						}
						// EBR A ESCRIBIR
						EBR := globals.EBR{}
						// LEO EL PRIMER EBR
						tempEbr := globals.EBR{} // VARIABLE PARA RECORRER LA LISTA DE EBR
						size := int(unsafe.Sizeof(tempEbr))
						file.Seek(int64(globals.ByteToInt(extendedPartition.Part_start[:])), 0)
						data := globals.ReadBytes(file, size)
						bufferebr := bytes.NewBuffer(data)
						err1 := binary.Read(bufferebr, binary.BigEndian, &tempEbr)
						if err1 != nil {
							log.Fatal("Error ", err1)
							return
						}
						//VARIABLE PARA GUARDAR DONDE SE VA A ESCRIBIR EL EBR
						startToWrite := 0
						if globals.ByteToString(tempEbr.Part_start[:]) == "-1" {
							// CALCULA QUE LA PARTICION QUEPA
							if (globals.ByteToInt(extendedPartition.Part_size[:]) - int(unsafe.Sizeof(EBR))) >= (cmd.Size * multiplicator) {
								// CALCULA DONDE INICIA LA PRIMERA PARTICION LOGICA
								start := globals.ByteToInt(extendedPartition.Part_start[:]) + int(unsafe.Sizeof(EBR)) + 1
								startToWrite = globals.ByteToInt(extendedPartition.Part_start[:])
								// ASIGNA LOS VALORES AL PRIMER EBR
								copy(EBR.Part_status[:], []byte("0"))
								copy(EBR.Part_fit[:], []byte(fit))
								copy(EBR.Part_start[:], []byte(strconv.Itoa(start)))
								copy(EBR.Part_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator)))
								copy(EBR.Part_next[:], []byte("-1"))
								copy(EBR.Part_name[:], []byte(cmd.Name))
							} else {
								fmt.Println("Error: la particion lógica " + cmd.Name + " no cabe en el disco " + cmd.Path)
								return
							}
						} else {
							for globals.ByteToInt(tempEbr.Part_next[:]) != -1 {
								size := int(unsafe.Sizeof(tempEbr))
								file.Seek(int64(globals.ByteToInt(tempEbr.Part_next[:])), 0)
								data := globals.ReadBytes(file, size)
								buffer := bytes.NewBuffer(data)
								err1 := binary.Read(buffer, binary.BigEndian, &tempEbr)
								if err1 != nil {
									log.Fatal("Error ", err1)
									return
								}
								//fmt.Println(globals.ByteToString(tempEbr.Part_name[:]), globals.ByteToString(tempEbr.Part_status[:]))
							}
							if ((globals.ByteToInt(extendedPartition.Part_start[:]) + globals.ByteToInt(extendedPartition.Part_size[:])) - (globals.ByteToInt(tempEbr.Part_start[:]) + globals.ByteToInt(tempEbr.Part_size[:]))) >= (cmd.Size * multiplicator) {
								// CALCULA DONDE INICIA LA PARTICION LOGICA
								start := globals.ByteToInt(tempEbr.Part_start[:]) + globals.ByteToInt(tempEbr.Part_size[:]) + int(unsafe.Sizeof(EBR)) + 2
								// CALCULA DONDE SE VA A ESCRIBIR EL EBR
								startToWrite = globals.ByteToInt(tempEbr.Part_start[:]) + globals.ByteToInt(tempEbr.Part_size[:]) + 1
								// ASIGNA LOS VALORES AL PRIMER EBR
								copy(EBR.Part_status[:], []byte("0"))
								copy(EBR.Part_fit[:], []byte(fit))
								copy(EBR.Part_start[:], []byte(strconv.Itoa(start)))
								copy(EBR.Part_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator)))
								copy(EBR.Part_next[:], []byte("-1"))
								copy(EBR.Part_name[:], []byte(cmd.Name))
								// APUNTO EL SIGUIENTE DE TEMPORAL AL INICIO DEL ACTUAL CREADO
								copy(tempEbr.Part_next[:], []byte(strconv.Itoa(startToWrite)))
								// REESCRIBO EL EBR TEMPORAL
								//fmt.Println(cmd.Name, startToWrite)
								//fmt.Println(globals.ByteToString(tempEbr.Part_name[:]), globals.ByteToInt(tempEbr.Part_start[:])-(int(unsafe.Sizeof(tempEbr))+1))
								file.Seek(int64(globals.ByteToInt(tempEbr.Part_start[:])-(int(unsafe.Sizeof(tempEbr))+1)), 0)
								var bufferControlTemp bytes.Buffer
								binary.Write(&bufferControlTemp, binary.BigEndian, &tempEbr)
								globals.WriteBytes(file, bufferControlTemp.Bytes())
							} else {
								fmt.Println("Error: la particion lógica " + cmd.Name + " no cabe en el disco " + cmd.Path)
								return
							}
						}
						// ESCRIBO EL EBR EN EL DISCO
						file.Seek(int64(startToWrite), 0)
						var bufferControl bytes.Buffer
						binary.Write(&bufferControl, binary.BigEndian, &EBR)
						globals.WriteBytes(file, bufferControl.Bytes())
						// IMPRIMO LOS DATOS DEL EBR CREADO
						globals.PrintEBR(EBR)
					} else {
						fmt.Println("Error: la particion " + cmd.Name + " no puede ser creada debido no existe particion extendida")
						return
					}
				}
				// REESCRIBO EL MBR EN EL DISCO
				file.Seek(0, 0)
				var bufferControl bytes.Buffer
				binary.Write(&bufferControl, binary.BigEndian, &MBR)
				globals.WriteBytes(file, bufferControl.Bytes())

				/*for i := 0; i < 4; i++ {
					fmt.Println(globals.ByteToString(MBR.Partitions[i].Part_start[:]))
				}*/
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
