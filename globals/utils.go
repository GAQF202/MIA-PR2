package globals

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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
	return string(byteArray[:])
}
