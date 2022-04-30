package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func Pause() {

	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)
	fmt.Println("Pause: Presiona cualquier letra para continuar[*]")
	for {
		stdin, _ := <-ch
		// AL DETECTAR CUALQUIER ENTRADA DE TECLADO ROMPE EL CICLO
		if stdin != "" {
			fmt.Println("\nSaliendo del Pause...")
			break
		}

	}
}
