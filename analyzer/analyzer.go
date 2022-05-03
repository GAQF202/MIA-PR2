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

func readFile(script_path string) string {
	// LEER EL ARREGLO DE BYTES DEL ARCHIVO
	datosComoBytes, err := ioutil.ReadFile(script_path /*"./test.txt"*/)
	if err != nil {
		log.Fatal(err)
	}
	// CONVERTIR EL ARREGLO A STRING
	datosComoString := string(datosComoBytes)

	return datosComoString
}

func AnalyzerF(script_path string, isFile bool) []globals.Command {
	var tree = make([]globals.Command, 0)
	input := ""
	// VERIFICO SI LA ENTRADA ES PARA EJECUTAR UN SCRIPT O DIRECTAMENTE UN COMANDO
	if isFile {
		// LEO EL ARCHIVO DE ENTRADA
		input = strings.ToLower(readFile(script_path)) + "\n"
	} else {
		// ASIGNO A LA ENTRADA EL COMANDO
		input = script_path
	}
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

	// VARIABLE PARA MANEJO DE COMENTARIOS
	isComment := false

	// BANDERA QUE INDICA QUE AUN NO ENCUENTRA UN CARACTER DISTINTO DE " " DESPUES DE "="
	valueFound := false
	catchSpaces := false

	for i, character := range input {
		letter := string(character)
		// BANDERA PARA SABER CUANDO INICIA UN COMENTARIO
		if letter == "#" {
			tree = append(tree, tempCommand)
			tempCommand = newCommand("comment")
			isComment = true
		}

		if !isComment {
			if !findValue {
				if letter != " " && letter != "\n" {
					tempWord += letter
					// COMANDOS PARA DISCOS
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
					} else if tempWord == "mount" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkfs" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
						// COMANDOS PARA SISTEMAS DE ARCHIVOS
					} else if tempWord == "mkdir" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkfile" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rep" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
						// COMANDOS PARA USUARIOS
					} else if tempWord == "login" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "logout" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkgrp" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rmgrp" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "mkuser" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "rmusr" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "exec" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
					} else if tempWord == "pause" {
						tree = append(tree, tempCommand)
						tempCommand = newCommand(tempWord)
						tempWord = ""
						// PARAMETROS
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
					} else if tempWord == "-ruta=" {
						tempPar = newParameter("ruta", "", -1)
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
					} else if tempWord == "-id=" {
						tempPar = newParameter("id", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-cont=" {
						tempPar = newParameter("cont", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-usuario=" {
						tempPar = newParameter("usuario", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-password=" {
						tempPar = newParameter("password", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-pwd=" {
						tempPar = newParameter("pwd", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
					} else if tempWord == "-grp=" {
						tempPar = newParameter("grp", "", -1)
						tempWord = ""
						findValue = true
						isIntValue = false
						// PARAMETROS DE UNA SOLA LETRA
					} else if tempWord == "-p" {
						// SI LLEGO AL FINAL DE LA CADENA ENTONCES GUARDA EL PARAMETRO Y EL COMANDO
						if i == (len(input) - 1) {
							tempPar = newParameter("-p", "-p", -1)
							tempWord = ""
							tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
							// GUARDO COMANDO EN EL ARBOL DE COMANDOS
							tree = append(tree, tempCommand)
						} else {
							// SI EL SIGUIENTE ES UN ESPACIO O UN "-" GUARDA EL PARAMETRO
							if string(input[i+1]) == " " || string(input[i+1]) == "-" || string(input[i+1]) == "\n" || string(input[i+1]) == "#" {
								tempPar = newParameter("-p", "-p", -1)
								tempWord = ""
								tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
								// SI NO ES UN CARACTER DE SEPARACION CONTINUA ANALIZANDO
							} else {
								continue
							}
						}
					} else if tempWord == "-r" {
						// SI LLEGO AL FINAL DE LA CADENA ENTONCES GUARDA EL PARAMETRO Y EL COMANDO
						if i == (len(input) - 1) {
							tempPar = newParameter("-r", "-r", -1)
							tempWord = ""
							tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
							// GUARDO COMANDO EN EL ARBOL DE COMANDOS
							tree = append(tree, tempCommand)
						} else {
							// SI EL SIGUIENTE ES UN ESPACIO O UN "-" GUARDA EL PARAMETRO
							if string(input[i+1]) == " " || string(input[i+1]) == "-" || string(input[i+1]) == "\n" || string(input[i+1]) == "#" {
								tempPar = newParameter("-r", "-r", -1)
								tempWord = ""
								tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
								// SI NO ES UN CARACTER DE SEPARACION CONTINUA ANALIZANDO
							} else {
								continue
							}
						}
					}
				}
				// SI ES UN COMENTARIO
				/*} else if isComment {
					tempWord += letter
					fmt.Println(tempWord)
					if letter == "\n" {
						fmt.Println("Entraa")
						fmt.Println(tempWord)
						tempWord = ""
						isComment = false
						findValue = false
					}
					// SI NO ES COMANDO NI COMENTARIO ES VALOR DE UN PARAMETRO
				} */
			} else {
				tempWord += letter
				// ACTIVAR BANDER SI VIENE UNA CADENA CON COMILLAS
				if letter == "\"" {
					if catchSpaces {
						catchSpaces = false
					} else {
						catchSpaces = true
					}
				}
				// CAPUTURO ESPACIOS
				if letter == " " && catchSpaces && valueFound {
					tempWord += letter
				}
				if letter == " " && !valueFound {
					continue
				} else {
					valueFound = true
				}
				// SI ENCUENTRA UN ESPACIO QUIERE DECIR SIGUE BUSCANDO COMANDOS
				if (letter == " " || i == (len(input)-1) || letter == "\n") && valueFound && !catchSpaces {
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
		} else {
			tempWord += letter
			// MIENTRAS EL CARACTER DE LA CADENA NO SEA UN SALTO
			// O NO SE HAYA LLEGADO AL FINAL DE LA CADENA SE SIGUE CONCATENDANDO
			if letter == "\n" || i == (len(input)-1) {
				tempPar = newParameter("value", strings.TrimSpace(tempWord), -1)
				tempCommand.Parameters = append(tempCommand.Parameters, tempPar)
				tempWord = ""
				isComment = false
				// SI LLEGA AL FINAL GUARDA EL COMANDO EN EL ARBOL DE COMANDOS
				if i == (len(input) - 1) {
					tree = append(tree, tempCommand)
				}
			}
		}

		// SI LLEGA AL FINAL ALMACENA EL ULTIMO COMANDO AL ARBOLITO
		if i == (len(input) - 1) {
			tree = append(tree, tempCommand)
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
