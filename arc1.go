// servidor com criacao dinamica de thread de servico
//
// go run arc1.go

package main

import (
	"fmt"
	"math/rand"
)

const (
        MAX_TH = 100
	NCL = 10
)

type Request struct {
	v int
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, " resp:", r)
	}
}

// ------------------------------------
// servidor
// thread de servico
func trataReq(id int, req Request) {
	fmt.Println(" trataReq ", id)
	req.ch_ret <- req.v * 2
}

// servidor que dispara threads de servico
func servidor(in chan Request, fim chan int) {
	var j int = 0
	for {
		j++
		req := <-in
		go trataReq(j, req)
		if j > MAX_TH {
			fim <- 1
		}
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("Criacao dinamica de threads")
	serv_chan := make(chan Request)
	fim_chan := make(chan int)
	go servidor(serv_chan, fim_chan)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}
	<-fim_chan
}