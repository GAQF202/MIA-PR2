package commands

import (
	"fmt"
	"os"

	"github.com/PR2_MIA/globals"
)

type RmdiskCmd struct {
	Path string
}

func (cmd *RmdiskCmd) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		}
	}
}

func (cmd *RmdiskCmd) Rmdisk() {
	err := os.Remove(cmd.Path)
	if err != nil {
		fmt.Println("Error: al eliminar el disco en la ruta " + cmd.Path)
	}
}
