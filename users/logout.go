package users

import (
	"github.com/PR2_MIA/globals"
)

func Logout() {
	// CIERRO SESION
	globals.GlobalUser.Logged = -1
	globals.GlobalUser.Uid = ""
	globals.GlobalUser.User_name = ""
	globals.GlobalUser.Pwd = ""
	globals.GlobalUser.Grp = ""
	globals.GlobalUser.Id_partition = ""
	globals.GlobalUser.Gid = ""
}
