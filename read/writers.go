package read

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/PR2_MIA/globals"
)

// FUNCION PARA ESCRIBIR BLOQUES DE CARPETAS
func WriteFileBlocks(file *os.File, position int, fileBlock globals.FileBlock) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &fileBlock)
	globals.WriteBytes(file, bufferControlBlocks.Bytes())
}

// FUNCION PARA ESCRIBIR BLOQUES DE ARCHIVOS
func WriteArchiveBlocks(file *os.File, position int, archiveBlock globals.ArchiveBlock) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &archiveBlock)
	globals.WriteBytes(file, bufferControlBlocks.Bytes())
}

// FUNCION PARA ESCRIBIR TABLAS DE INODOS
func WriteInodes(file *os.File, position int, inode globals.InodeTable) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &inode)
	globals.WriteBytes(file, bufferControlBlocks.Bytes())

}

// FUNCION PARA ESCRIBIR EL SUPERBLOQUE EN EL ARCHIVO
func WriteSuperBlock(file *os.File, position int, super_bloque globals.SuperBloque) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &super_bloque)
	globals.WriteBytes(file, bufferControlBlocks.Bytes())
}

// FUNCION PARA ESCRIBIR BITMAPS
func WriteBitmap(file *os.File, position int, bitmap []byte) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &bitmap)
	globals.WriteBytes(file, bufferControlBlocks.Bytes())
}
