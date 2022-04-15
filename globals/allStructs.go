package globals

type Parameter struct {
	Name        string
	StringValue string
	IntValue    int
}

type Command struct {
	Parameters []Parameter
	Name       string
}

type VoidSpace struct {
	Start int
	Size  int
}

type Partition struct {
	Part_status [3]byte
	Part_type   [4]byte
	Part_fit    [4]byte
	Part_start  [16]byte
	Part_size   [16]byte
	Part_name   [16]byte
}

type MBR struct {
	Mbr_size           [16]byte
	Mbr_fecha_creacion [18]byte
	Mbr_dsk_signature  [16]byte
	Dsk_fit            [4]byte
	Partitions         [4]Partition
}

type EBR struct {
	Part_status [3]byte
	Part_fit    [4]byte
	Part_start  [16]byte
	Part_size   [16]byte
	Part_next   [16]byte
	Part_name   [16]byte
}
