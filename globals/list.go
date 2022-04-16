package globals

import (
	"fmt"
	"strconv"
)

func GetLetter(number int) string {
	letter := ""
	for ch := 97; ch <= 122; ch++ {
		if ch == (number + 96) {
			letter = string(ch)
		}
	}
	return letter
}

type PartitionMounted struct {
	DiskName      string // DISCO DONDE SE ENCUENTRA LA PARTICION MONTADA
	PartitionName string // NOMBRE DE LA PARTICION MONTADA
	Id            string // ID DE LA PARTICION MONTADA
	Start         int    // INICIO DE LA PARTICION MONTADA EN EL DISCO
	Status        int    // ATRIBUTO PARA SABER SI ESTA MONTADA LA PARTICION
	Type          string // TIPO DE PARTICION MONTADA
	Fit           string // TIPO DE FIT DE LA PARTICION MONTADA
	Size          int    // TAMANIO DE LA PARTICION MONTADA
	Path          string // RUTA DONDE SE ENCUENTRA EL DISCO DE LA PARTICION
}

type Element struct {
	DiskName          string
	DiskNumber        int
	Path_disk         string // RUTA DEL DISCO
	PartitionsMounted []PartitionMounted
}

type ListPartitionsMounted struct {
	Partitions []Element
}

// FUNCION PARA CREAR NUEVO ELEMENTO
func NewElement(disk_name string, disk_number int, path_disk string) Element {
	slice := make([]PartitionMounted, 1)
	return Element{disk_name, disk_number, path_disk, slice}
}

// FUNCION PARA CREAR UNA NUEVA LISTA DE PARTICIONES
func NewRam() ListPartitionsMounted {
	slice := make([]Element, 1)
	return ListPartitionsMounted{slice}
}

func (list *ListPartitionsMounted) MountPartition(disk_name string, partition_name string, start int, path string, type_partition string, fit string, size int) {
	idNumber := -1
	idLeter := ""

	// BUSCA SI YA SE HA MONTADO ALGUNA PARTICION DEL DISCO
	for i, partition := range list.Partitions {
		if partition.DiskName == disk_name {
			for _, par := range list.Partitions[i].PartitionsMounted {
				if par.PartitionName == partition_name && par.Status == 1 {
					fmt.Println("Error: la particion " + partition_name + " ya está montada en ram " + disk_name)
					return
				}
			}
			idNumber = partition.DiskNumber
			idLeter = GetLetter(len(partition.PartitionsMounted) + 1)
			id := "57" + strconv.Itoa(idNumber) + idLeter
			newPartition := PartitionMounted{disk_name, partition_name, id, start, 1, type_partition, fit, size, path}
			list.Partitions[i].PartitionsMounted = append(list.Partitions[i].PartitionsMounted, newPartition)
			break
		}
	}

	// SI NO SE ENCONTRO EL DISCO LO CREO
	if idNumber == -1 {
		idNumber = len(list.Partitions)
		idLeter = "a"
		id := "57" + strconv.Itoa(idNumber) + idLeter
		element := NewElement(disk_name, len(list.Partitions), path)
		list.Partitions = append(list.Partitions, element)
		newPartition := PartitionMounted{disk_name, partition_name, id, start, 1, type_partition, fit, size, path}
		list.Partitions[len(list.Partitions)-1].PartitionsMounted = append(list.Partitions[len(list.Partitions)-1].PartitionsMounted, newPartition)
	}

	//fmt.Println(idNumber, idLeter, disk_name)
}

func (list *ListPartitionsMounted) GetElement(myId string) PartitionMounted {
	res := PartitionMounted{}
	// BUSCA SI YA SE HA MONTADO ALGUNA PARTICION DEL DISCO
	for _, partition := range list.Partitions {
		for _, par := range partition.PartitionsMounted { //list.Partitions[i].PartitionsMounted {
			if par.Id == myId {
				res = par
				break
			}
		}
	}
	return res
}

var GlobalList ListPartitionsMounted = NewRam()
