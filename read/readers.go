package read

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

// FUNCION PARA LEER BITMPAS DEL DISCO
func ReadBitMap(file *os.File, position int, bitmap_size int) []byte {
	var bitmap = make([]byte, bitmap_size)

	size := int(bitmap_size)
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &bitmap)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return bitmap
}

// FUNCION PARA LEER SUPERBLOQUE
func ReadSuperBlock(file *os.File, position int) globals.SuperBloque {
	var super_bloque = globals.SuperBloque{}

	size := int(unsafe.Sizeof(super_bloque))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &super_bloque)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return super_bloque
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadInode(file *os.File, position int) globals.InodeTable {
	var inode = globals.InodeTable{}

	size := int(unsafe.Sizeof(inode))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &inode)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return inode
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadFileBlock(file *os.File, position int) globals.FileBlock {
	var block = globals.FileBlock{}

	size := int(unsafe.Sizeof(block))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &block)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return block
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadArchiveBlock(file *os.File, position int) globals.ArchiveBlock {
	var block = globals.ArchiveBlock{}

	size := int(unsafe.Sizeof(block))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &block)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return block
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadFileMbr(file *os.File, position int) globals.MBR {
	var mbr = globals.MBR{}

	size := int(unsafe.Sizeof(mbr))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &mbr)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return mbr
}

// FUNCION PARA LEER PARTICIONES DEL DISCO
func ReadPartition(file *os.File, position int) globals.Partition {
	var partition = globals.Partition{}

	size := int(unsafe.Sizeof(partition))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &partition)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return partition
}

// FUNCION PARA LEER EBR DEL DISCO
func ReadEbr(file *os.File, position int) globals.EBR {
	var ebr = globals.EBR{}

	size := int(unsafe.Sizeof(ebr))
	file.Seek(int64(position), 0)
	data := globals.ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &ebr)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return ebr
}

// FUNCION PARA BUSCAR EL ULTIMO INODO DE UNA RUTA
func GetInodeWithPath(path string, real_path string, start int) globals.InodeTable {
	//fmt.Println("entraaa en el file")
	if globals.GlobalUser.Logged == -1 {
		fmt.Println("Error: Para utilizar mkfile necesitas estar logueado")
		return globals.InodeTable{}
	}

	// ABRO EL ARCHIVO
	file, err := os.OpenFile(real_path, os.O_RDWR, 0777)
	// VERIFICACION DE ERROR AL ABRIR EL ARCHIVO
	if err != nil {
		log.Fatal("Error ", err)
		return globals.InodeTable{}
	}

	// OBTENGO TODAS LAS CARPETAS PADRES ANTES DEL ARCHIVO
	parent_path := path[:strings.LastIndex(path, "/")]

	//current_date := globals.GetDate()

	// OBTENGO TODAS LAS CARPETAS PADRES DEL ARCHIVO
	routes := strings.Split(parent_path, "/")

	//saltar_busqueda := false
	exist_route := false

	// SI EL TAMANIO DE LAS RUTAS SEPARADAS POR / ES
	// CERO QUIERE DECIR QUE EL ARCHIVO SE DEBE CREAR EN LA RAIZ
	if len(routes) == 0 && path == "/" {
		exist_route = true
	} else {
		temp := []string{"/"}
		temp = append(temp, routes[1:]...)
		routes = temp
	}

	// LEO EL SUPERBLOQUE
	super_bloque := globals.SuperBloque{}
	super_bloque = ReadSuperBlock(file, start)

	// CREACION DE ARRAY PARA ALMACENAR LOS BITMPAS
	var bitinodes = make([]byte, globals.ByteToInt(super_bloque.Inodes_count[:]))
	var bitblocks = make([]byte, globals.ByteToInt(super_bloque.Blocks_count[:]))
	bitinodes = ReadBitMap(file, globals.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
	bitblocks = ReadBitMap(file, globals.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

	// VERIFICA QUE EXISTAN LAS CARPETAS ANTES DEL ARCHIVO

	temp_inode := globals.InodeTable{}

	//LEO EL PRIMER INODO
	temp_inode = ReadInode(file, globals.ByteToInt(super_bloque.Inode_start[:]))

	// VECTOR PARA GUARDAR LAS RUTAS QUE FALTAN POR CREARSE
	var remaining_routes = make([]string, len(routes))
	// CREO UNA COPIA PARA QUE NO SE ALTERE EL ROUTE
	copy(remaining_routes, routes)

	// RECORRE LA RUTA
	for path_index := 0; path_index < len(routes); path_index++ {
		exist_path := false
		// RECORRE LOS PUNTEROS DEL INODO
		for pointerIndex := 0; pointerIndex < 16; pointerIndex++ {
			// RECORRO SOLO LOS BLOQUES DE LOS INODOS DE TIPO CARPETA
			if temp_inode.Block[pointerIndex] != -1 && globals.ByteToString(temp_inode.Type[:]) == "0" {
				file_block := globals.FileBlock{}
				file_block = ReadFileBlock(file, (globals.ByteToInt(super_bloque.Block_start[:]) + (int(temp_inode.Block[pointerIndex]) * int(unsafe.Sizeof(file_block)))))
				// RECORRE LOS PUNTEROS DE LOS BLOQUES
				for blockIndex := 0; blockIndex < 4; blockIndex++ {
					if file_block.Content[blockIndex].Inodo != -1 {
						if globals.ByteToString(file_block.Content[blockIndex].Name[:]) == routes[path_index] {
							// ELIMINO LAS RUTAS QUE YA ESTAN CREADAS PARA QUE QUEDEN SOLO LAS RESTANTES
							if len(remaining_routes) == len(routes) {
								remaining_routes = globals.RemoveIndex(remaining_routes, 0)
								remaining_routes = globals.RemoveIndex(remaining_routes, 0)
							} else {
								remaining_routes = globals.RemoveIndex(remaining_routes, 0)
							}
							temp_inode = ReadInode(file, globals.ByteToInt(super_bloque.Inode_start[:])+(int(file_block.Content[blockIndex].Inodo)*int(unsafe.Sizeof(temp_inode))))
							exist_route = true
							exist_path = true
						}
					}
				}
			}
		}
		if !exist_path {
			exist_route = false
		}
	}
	// VALIDACION PARA SABER SI EL ARCHIVO SE CREA EN LA RAIZ
	if routes[0] == "/" && len(routes) == 1 {
		temp_inode = ReadInode(file, globals.ByteToInt(super_bloque.Inode_start[:]))
		exist_route = true
	}
	//VERIFICACION DE EXISTENCIA DE RUTAS
	if !exist_route {
		fmt.Println("No existe la ruta indicada")
		return globals.InodeTable{}
	}
	return temp_inode
}
