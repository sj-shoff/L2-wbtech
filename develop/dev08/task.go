package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

// GetTime получает текущее время с NTP-сервера "0.pool.ntp.org"
// и выводит его в стандартный вывод. В случае ошибки выводит сообщение
// об ошибке в стандартный поток ошибок.
func GetTime() {
	time, err := ntp.Time("0.pool.ntp.org")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	} else {
		fmt.Println(time)
	}
}

func main() {
	GetTime()
}
