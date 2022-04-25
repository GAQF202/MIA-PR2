package read

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
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
