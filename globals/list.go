package globals

import "fmt"

type DisksMounted struct {
	DiskName   string
	DiskNumber int
}

type PartitionMounted struct {
	DiskName      string // DISCO DONDE SE ENCUENTRA LA PARTICION MONTADA
	PartitionName string // NOMBRE DE LA PARTICION MONTADA
	Id            string // ID DE LA PARTICION MONTADA
	Start         int    // INICIO DE LA PARTICION MONTADA EN EL DISCO
	Path_disk     string // RUTA DEL DISCO
	Status        bool   // ATRIBUTO PARA SABER SI ESTA MONTADA LA PARTICION
}

type ListPartitionsMounted struct {
	DisksMounted      []DisksMounted
	PartitionsMounted []PartitionMounted
}

func (list *ListPartitionsMounted) MountPartition(disk_name string, partition_name string, start int, path string) {
	idNumber := -1
	idLeter := ""
	// BUSCO SI YA SE HA MONTADO EL DISCO ANTERIORMENTE
	for disk_i, disk := range list.DisksMounted {
		if disk.DiskName == partition_name {
			idNumber = disk_i
		}
	}

	fmt.Println(idNumber, idLeter)
}

var GlobalList ListPartitionsMounted = ListPartitionsMounted{}
