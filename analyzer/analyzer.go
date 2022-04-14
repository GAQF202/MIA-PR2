package analyzer

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/PR2_MIA/globals"
)

func newParameter(Name string, StringValue string, IntValue int) globals.Parameter {
	return globals.Parameter{Name, StringValue, IntValue}
}

func newCommand(Name string) globals.Command {
	temp := make([]globals.Parameter, 1)
	return globals.Command{temp, Name}
}

func readFile() string {
	// LEER EL ARREGLO DE BYTES DEL ARCHIVO
	datosComoBytes, err := ioutil.ReadFile("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	// CONVERTIR EL ARREGLO A STRING
	datosComoString := string(datosComoBytes)

	return datosComoString
}

func AnalyzerF() []globals.Command {
	var tree = make([]globals.Command, 0)
	// LEO EL ARCHIVO DE ENTRADA
	input := strings.ToLower(readFile())
	//VARIABLE PARA GUARDAR LOS COMANDOS TEMPORALES
	tempCommand := newCommand("")
	tempPar := newParameter("", "", -1)
	//VARIABLE PARA RECONOCER COMANDOS
	tempWord := ""
	//VARIABLE PARA INDICAR QUE DEBE BUSCAR UN VALOR
	findValue := false

	tempIntValue := -1
	tempStringValue := ""
	isIntValue := false

	// BANDERA QUE INDICA QUE AUN NO ENCUENTRA UN CARACTER DISTINTO DE " " DESPUES DE "="
	valueFound := false

	for i, character := range input {
		letter := string(character)

		if !findValue {
			if letter != " " && letter != "\n" {
				tempWord += letter
				if tempWord == "mkdisk" {
					tree = append(tree, tempCommand)
					tempCommand = newCommand(tempWord)
					tempWord = ""
				} else if tempWord == "rmdisk" {
					tree = append(tree, tempCommand)
					tempCommand = newCommand(tempWord)
					tempWord = ""
				} else if tempWord == "fdisk" {
					tree = append(tree, tempCommand)
					tempCommand = newCommand(tempWord)
					tempWord = ""
				} else if tempWord == "-size=" {
					tempPar = newParameter("size", "", -1)
					tempWord = ""
					findValue = true
					isIntValue = true
				} else if tempWord == "-fit=" {
					tempPar = newParameter("fit", "", -1)
					tempWord = ""
					findValue = true
					isIntValue = false
				} else if tempWord == "-unit=" {
					tempPar = newParameter("unit", "", -1)
					tempWord = ""
					findValue = true
					isIntValue = false
				} else if tempWord == "-path=" {
					tempPar = newParameter("path", "", -1)
					tempWord = ""
					findValue = true
					isIntValue = false
				} else if tempWord == "-type=" {
					tempPar = newParameter("type", "", -1)
					tempWord = ""
					findValue = true
					isIntValue = false
				} else if tempWord == "-name=" {
					tempPar = newParameter("name", "", -1)
					tempWord = ""
					findValue = true
					isIntValue = false
				}
			}
		} else {
			tempWord += letter
			if letter == " " && !valueFound {
				continue
			} else {
				valueFound = true
			}
			// SI ENCUENTRA UN ESPACIO QUIERE DECIR SIGUE BUSCANDO COMANDOS
			if (letter == " " || i == (len(input)-1) || letter == "\n") && valueFound {
				if isIntValue {
					convertedValue, _ := strconv.Atoi(strings.TrimSpace(tempWord))
					tempIntValue = convertedValue
				} else {
					// DEVUELVE EL STRING SIN SALTOS
					tempStringValue = strings.TrimSuffix(strings.TrimSpace(tempWord), "\n")
				}
				tempPar.IntValue = tempIntValue
				tempPar.StringValue = tempStringValue
				tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
				tempWord = ""
				tempStringValue = ""
				tempIntValue = -1
				findValue = false
				valueFound = false // REGRESO BANDERA DE CARACTER DIFERENTE DE " "
				if i == (len(input) - 1) {
					tree = append(tree, tempCommand)
				}
			}
		}
	}

	/*fmt.Println(tempIntValue, tempStringValue, len(tree))
	for s, ver := range tree {
		if ver.Name != "" {
			fmt.Println(ver.Name)
			for _, parameter := range ver.Parameters {
				fmt.Println(parameter, s)
			}
		}
	}*/
	return tree
}
