package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"time"
)

var totalPorts = 0

func WriteLog(pts []int, fname, address string) {

	file, err := os.Create(fname)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(file, address)
	for i := range pts {
		fmt.Fprintf(file, "Порт %d открыт!\n", pts[i])
	}
}
func GenName(ad string) string {
	dt := time.Now()
	var a string = dt.Format("01-02-2006")
	return ad + "_" + a + ".txt"
}
func worker(ports chan int, address string, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", address, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
func main() {
	fmt.Println("Введите адрес сайта для сканирования: ")
	var address string

	fmt.Fscan(os.Stdin, &address)

	fmt.Println("Введите диапазон портов для сканирования: ")
	var lb, ub int
	fmt.Fscan(os.Stdin, &lb)
	fmt.Fscan(os.Stdin, &ub)
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	//Цикл для создания воркеров
	for i := 0; i < cap(ports); i++ {
		go worker(ports, address, results)
	}
	go func() {
		for i := lb; i <= ub; i++ {
			ports <- i
		}
	}()
	for i := lb; i < ub; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	sort.Ints(openports)
	fmt.Printf("\nВсего портов: %d\n--------- \n", len(openports))
	z := 1
	for _, port := range openports {
		fmt.Printf("%d) Порт %d открыт!\n", z, port)
		z++
	}

	fmt.Printf("Порты отсканированы. Желаете ли вы создать лог-файл \"%s\"?   (y|n) \n", GenName(address))

	var exit string
	fmt.Scan(&exit)
	if exit == "y" || exit == "Y" {
		WriteLog(openports, GenName(address), address)

	}
	close(ports)
	close(results)
}
