package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PR2_MIA/globals"
	"github.com/PR2_MIA/read"
)

type RepCmd struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func (cmd *RepCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "name" {
			cmd.Name = parameter.StringValue
		} else if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "id" {
			cmd.Id = parameter.StringValue
		} else if parameter.Name == "ruta" {
			cmd.Ruta = parameter.StringValue
		}
	}
}

func (cmd *RepCmd) Rep() {
	if cmd.Name == "" || cmd.Id == "" || cmd.Path == "" {
		fmt.Println("Error: faltan parametros obligatorios en el comando rep")
		return
	}

	// CREA LAS CARPETAS PADRE
	parent_path := cmd.Path[:strings.LastIndex(cmd.Path, "/")]
	if err := os.MkdirAll(parent_path, 0700); err != nil {
		log.Fatal(err)
	}

	// OBTENGO EL NOMBRE DEL REPORTE A GENERAR
	report_name := cmd.Path[strings.LastIndex(cmd.Path, "/")+1 : len(cmd.Path)]

	// OBTENGO LA PARTICION MONTADA
	partition_m := globals.GlobalList.GetElement(cmd.Id)
	// ABRO EL ARCHIVO
	file, err := os.OpenFile(partition_m.Path, os.O_RDWR, 0777)
	// VERIFICACION DE ERROR AL ABRIR EL ARCHIVO
	if err != nil {
		log.Fatal("Error ", err)
		return
	}

	// LEO EL MBR
	mbr := globals.MBR{}
	mbr = read.ReadFileMbr(file, 0)

	// LEO TODO LO RELACIONADO AL SISTEMA DE ARCHIVOS
	super_bloque := globals.SuperBloque{}
	super_bloque = read.ReadSuperBlock(file, partition_m.Start)

	// CREACION DE ARRAY PARA ALMACENAR LOS BITMPAS
	var bitinodes = make([]byte, globals.ByteToInt(super_bloque.Inodes_count[:]))
	var bitblocks = make([]byte, globals.ByteToInt(super_bloque.Blocks_count[:]))
	bitinodes = read.ReadBitMap(file, globals.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
	bitblocks = read.ReadBitMap(file, globals.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitblocks))

	if cmd.Name == "disk" {
		var porcentage float64 = 0

		dotContent := `digraph html { abc [shape=none, margin=0, label=< 
			<TABLE BORDER="1" COLOR="#10a20e" CELLBORDER="1" CELLSPACING="3" CELLPADDING="4">`

		logicas := "\n<TR>"
		all_partitions := "\n<TR>\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\">MBR</TD>\n"

		avalaible_space := 0

		for i := 0; i < 4; i++ {
			if globals.ByteToString(mbr.Partitions[i].Part_type[:]) != "p" {

				colspan := 2
				temp := globals.EBR{} // GUARDA EL TEMPORAL PARA RECORRER LA LISTA

				// LEO LA PRIMERA PARTICION LOGICA A DONDE APUNTA LA EXTENDIDA Y ASIGNO A TEMP
				temp = read.ReadEbr(file, globals.ByteToInt(mbr.Partitions[i].Part_start[:]))
				// GRAFICA DE LOGICA
				porcentage = float64(float64((globals.ByteToInt(temp.Part_size[:]))*100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))
				logicas += "\n<TD COLOR=\"#87b8a4\">EBR</TD>\n"
				logicas += "\n<TD COLOR=\"#87b8a4\"> Lógica <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"

				// MIENTRAS NO LLEGUE AL FINAL DE LA LISTA
				for globals.ByteToInt(temp.Part_next[:]) != -1 {
					colspan += 2
					temp = read.ReadEbr(file, globals.ByteToInt(temp.Part_next[:]))
					// GRAFICA DE LOGICA
					porcentage = float64((float64(globals.ByteToInt(temp.Part_size[:])) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))
					logicas += "\n<TD COLOR=\"#87b8a4\">EBR</TD>\n"
					logicas += "\n<TD COLOR=\"#87b8a4\"> Lógica <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
				}

				// MIENTRAS NO LLEGUE A LA ULTIMA PARTICION
				if i != 3 {
					// GRAFICA DE EXTENDIDA
					porcentage = float64((float64(globals.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))
					all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Extendida <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"

					if globals.ByteToInt(mbr.Partitions[i+1].Part_size[:]) != -1 {
						// CALCULO ESPACIO VACIO
						avalaible_space = globals.ByteToInt(mbr.Partitions[i+1].Part_start[:]) - (globals.ByteToInt(mbr.Partitions[i].Part_start[:]) + globals.ByteToInt(mbr.Partitions[i].Part_size[:]))
						// CALCULO PORCENTAJE
						porcentage = float64((float64(avalaible_space))*100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:]))
						if porcentage > 0.8 {
							all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
						}
					}
				} else {
					// GRAFICA DE EXTENDIDA
					porcentage = float64((float64(globals.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))
					all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Extendida <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
					//all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Extendida <BR/>" + fmt.Sprint(porcentage) + "%</TD>\n"
					// CALCULO ESPACIO VACIO
					avalaible_space = globals.ByteToInt(mbr.Mbr_size[:]) - (globals.ByteToInt(mbr.Partitions[i].Part_start[:]) + globals.ByteToInt(mbr.Partitions[i].Part_size[:]))
					// CALCULO PORCENTAJE
					porcentage = float64((float64(avalaible_space) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))
					if porcentage > 0.8 {
						//all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Libre <BR/>" + fmt.Sprint(porcentage) + "%</TD>\n"
						all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
					}
				}
			} else {
				// GRAFICO SOLO LAS PARTICIONES EXISTENTES
				if globals.ByteToInt(mbr.Partitions[i].Part_size[:]) != -1 {
					// MIENTRAS NO LLEGUE A LA ULTIMA PARTICION
					if i != 3 {
						// GRAFICA DE PRIMARIA
						porcentage = float64((float64(globals.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))

						all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Primaria <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"

						if globals.ByteToInt(mbr.Partitions[i+1].Part_size[:]) != -1 {
							// CALCULO ESPACIO VACIO
							avalaible_space = globals.ByteToInt(mbr.Partitions[i+1].Part_start[:]) - (globals.ByteToInt(mbr.Partitions[i].Part_start[:]) + globals.ByteToInt(mbr.Partitions[i].Part_size[:]))
							// CALCULO PORCENTAJE
							porcentage = float64((float64(avalaible_space) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))

							if porcentage > 0.8 {
								all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
							}
						} else {
							// CALCULO ESPACIO VACIO
							avalaible_space = globals.ByteToInt(mbr.Mbr_size[:]) - (globals.ByteToInt(mbr.Partitions[i].Part_start[:]) + globals.ByteToInt(mbr.Partitions[i].Part_size[:]))
							// CALCULO PORCENTAJE
							porcentage = float64((float64(avalaible_space) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))

							if porcentage > 0.8 {
								all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
							}
						}
					} else {
						// GRAFICA DE PRIMARIA
						porcentage = float64((float64(globals.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(globals.ByteToInt(mbr.Mbr_size[:])))
						all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Primaria <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
						// CALCULO ESPACIO VACIO
						avalaible_space = globals.ByteToInt(mbr.Mbr_size[:]) - (globals.ByteToInt(mbr.Partitions[i].Part_start[:]) + globals.ByteToInt(mbr.Partitions[i].Part_size[:]))
						// CALCULO PORCENTAJE
						porcentage = float64((avalaible_space * 100) / globals.ByteToInt(mbr.Mbr_size[:]))

						if porcentage > 0.8 {
							all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprint(porcentage) + "%</TD>\n"
						}
					}
				}
			}
		}

		all_partitions += "</TR>\n"
		logicas += "</TR>\n"
		dotContent += all_partitions + logicas + "</TABLE>>];\n}"

		fmt.Println(dotContent)
		fmt.Println(report_name)
	}
}
