package commands

import (
	"github.com/PR2_MIA/globals"
)

type Comment struct {
	Value string
}

func (cmd *Comment) AssignParameters(command globals.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "value" {
			cmd.Value = parameter.StringValue
		}
	}
}

func (cmd *Comment) ShowComment() {
	//fmt.Println(cmd.Value)
}
