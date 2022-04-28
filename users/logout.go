package users

import (
	"fmt"

	"github.com/PR2_MIA/globals"
)

func Logout() {
	if globals.GlobalUser.Logged != -1 {
		// CIERRO SESION
		globals.GlobalUser.Logged = -1
		globals.GlobalUser.Uid = ""
		globals.GlobalUser.User_name = ""
		globals.GlobalUser.Pwd = ""
		globals.GlobalUser.Grp = ""
		globals.GlobalUser.Id_partition = ""
		globals.GlobalUser.Gid = ""
	} else {
		fmt.Println("Error: no se puede realizar el logout ya que no hay ningun usuario logueado actualmente")
	}
}
