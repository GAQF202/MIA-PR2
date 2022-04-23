package globals

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func WriteBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func GetDate() string {

	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return fecha
}

func GetRandom() int {
	res := rand.Intn(319)
	return res
}

func ByteToString(byteArray []byte) string {
	temp := bytes.Trim(byteArray, "\x00")
	str := string(temp[:])
	return str
}

func ByteToInt(byteArray []byte) int {
	return ToInt(ByteToString(byteArray))
}

func ToInt(stringValue string) int {
	res, err := strconv.Atoi(stringValue)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func BubbleSort(array []Partition) {
	for i := 0; i < len(array)-1; i++ {
		for j := 0; j < len(array)-i-1; j++ {
			if ToInt(ByteToString(array[j].Part_start[:])) > ToInt(ByteToString(array[j+1].Part_start[:])) {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}

func SortFreeSpaces(array []VoidSpace) {
	for i := 0; i < len(array)-1; i++ {
		for j := 0; j < len(array)-i-1; j++ {
			if array[j].Size > array[j+1].Size {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}

// FUNCION PARA BORRAR ELEMENTO POR INIDICE DE UN ARRAY
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
