package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/PR2_MIA/globals"
	"github.com/PR2_MIA/read"
)

func createPortTd(content string, port string) string {
	return "<td port=\"" + port + "\">" + content + "</td>"
}

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
	report_name := cmd.Path[:strings.LastIndex(cmd.Path, ".")+1]

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
	bitblocks = read.ReadBitMap(file, globals.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

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

		// CREO Y ESCRIBO EL ARCHIVO .dot
		err := ioutil.WriteFile(report_name+"dot", []byte(dotContent), 0644)
		if err != nil {
			log.Fatal(err)
		}
		// GRAFICO EL ARCHIVO .dot CREADO
		globals.GraphDot(report_name+"dot", cmd.Path)
	} else if cmd.Name == "tree" {
		// VARIABLE PARA RECORRER INODOS
		temp_inode := globals.InodeTable{}

		// VARIABLE PARA MOSTRAR TODOS LOS TIPOS DE BLOQUES
		file_block := globals.FileBlock{}
		archive_block := globals.ArchiveBlock{}

		nodes := ""
		blocks := ""
		edges := ""

		dotContent := "digraph {\ngraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\nnode [shape=plain]\nrankdir=LR;"

		// RECORRO INODOS
		for i := 0; i < len(bitinodes); i++ {
			// SI NO ES UN INODO LIBRE
			if bitinodes[i] != '0' {

				// LEO EL INODO
				temp_inode = read.ReadInode(file, globals.ByteToInt(super_bloque.Inode_start[:])+(i*int(unsafe.Sizeof(temp_inode))))
				//fmt.Println(globals.ByteToInt(super_bloque.Inode_start[:]) + (i * int(unsafe.Sizeof(temp_inode))))

				nodes += "inode" + strconv.Itoa(i) + " [label=< \n <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
				nodes += "<tr><td bgcolor=\"#01f5ab\">INODE</td><td bgcolor=\"#01f5ab\">" + strconv.Itoa(i) + "</td></tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("UID", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Uid[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("GID", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Gid[:]), "")
				nodes += "</tr>\n"

				nodes += "<tr>"
				nodes += createPortTd("SIZE", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Size[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("LECTURA", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Atime[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("CREACION", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Ctime[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("MODIFICACION", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Mtime[:]), "")
				nodes += "</tr>\n"

				for block_index := 0; block_index < 14; block_index++ {
					nodes += "<tr>"
					nodes += createPortTd("AP"+strconv.Itoa(block_index), "")
					nodes += createPortTd(strconv.Itoa(int(temp_inode.Block[block_index])), "i"+strconv.Itoa(i)+"b"+strconv.Itoa(int(temp_inode.Block[block_index])))
					nodes += "</tr>\n"

					if temp_inode.Block[block_index] != -1 {
						edges += "inode" + strconv.Itoa(i) + ":i" + strconv.Itoa(i) + "b" + strconv.Itoa(int(temp_inode.Block[block_index])) + "->" + "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + ";\n"
						// SI ES UN INODO DE ARCHIVO
						if globals.ByteToString(temp_inode.Type[:]) == "1" {
							archive_block = read.ReadArchiveBlock(file, globals.ByteToInt(super_bloque.Block_start[:])+(int(temp_inode.Block[block_index])*int(unsafe.Sizeof(archive_block))))

							// ENCABEZADO DEL BLOQUE
							blocks += "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + " [label=< \n <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
							blocks += "<tr><td bgcolor=\"#f6ec1e\">BLOCK</td><td bgcolor=\"#f6ec1e\">" + strconv.Itoa(int(temp_inode.Block[block_index])) + "</td></tr>\n"
							// ENCABEZADO EL CONTENIDO
							block_content := ""
							for con := 0; con < 64; con++ {
								block_content += string(archive_block.Content[con])
							}
							blocks += "<tr><td colspan=\"2\">" + block_content + "</td></tr>\n"
							// CIERRO LA TABLA
							blocks += "</table>>]; \n"
						} else if globals.ByteToString(temp_inode.Type[:]) == "0" {
							file_block = read.ReadFileBlock(file, globals.ByteToInt(super_bloque.Block_start[:])+(int(temp_inode.Block[block_index])*int(unsafe.Sizeof(file_block))))
							// ENCABEZADO DEL BLOQUE
							blocks += "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + " [label=< \n <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
							blocks += "<tr><td bgcolor=\"#f61e73\">BLOCK</td><td bgcolor=\"#f61e73\">" + strconv.Itoa(int(temp_inode.Block[block_index])) + "</td></tr>\n"

							for content := 0; content < 4; content++ {
								blocks += "<tr>" + createPortTd(globals.ByteToString(file_block.Content[content].Name[:]), "") + createPortTd(strconv.Itoa(int(file_block.Content[content].Inodo)), "b"+strconv.Itoa(int(temp_inode.Block[block_index]))+"i"+strconv.Itoa(int(file_block.Content[content].Inodo))) + "</tr>\n"
								if file_block.Content[content].Inodo != -1 && globals.ByteToString(file_block.Content[content].Name[:]) != "." && globals.ByteToString(file_block.Content[content].Name[:]) != ".." {
									edges += "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + ":b" + strconv.Itoa(int(temp_inode.Block[block_index])) + "i" + strconv.Itoa(int(file_block.Content[content].Inodo)) + "->" + "inode" + strconv.Itoa(int(file_block.Content[content].Inodo)) + ";\n"
								}
							}
							blocks += "</table>>]; \n"
						}
					}
				}
				nodes += "<tr>"
				nodes += createPortTd("TIPO", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Type[:]), "")
				nodes += "</tr>\n"

				nodes += "<tr>"
				nodes += createPortTd("PERMISOS", "")
				nodes += createPortTd(globals.ByteToString(temp_inode.Perm[:]), "")
				nodes += "</tr>\n"

				nodes += "</table>>]; \n"
			}
		}
		dotContent += nodes
		dotContent += blocks
		dotContent += edges
		dotContent += "\n}"

		//fmt.Println(dotContent)

		// CREO Y ESCRIBO EL ARCHIVO .dot
		err := ioutil.WriteFile(report_name+"dot", []byte(dotContent), 0644)
		if err != nil {
			log.Fatal(err)
		}
		// GRAFICO EL ARCHIVO .dot CREADO
		globals.GraphDot(report_name+"dot", cmd.Path)
	}
}
