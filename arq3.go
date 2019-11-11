// servidor com criacao dinamica de thread de servico
//
// go run arq1.go

package main

import (
	"fmt"
	"math/rand"
)

const (
     MAXTH = 10
	NCL = 1000
)

// estrutura que define o pedido dos clientes
type Request struct {
	// define o valor para o calculo da thread
	v int
	// define o canal para a resposta da thread
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request) {
	var v, r int

	// cria canal para receber resposta das threads
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)

		// envia o pedido para o servidor
		req <- Request{v, my_ch}

		// espera e le resposta de uma thread
		r = <-my_ch

		// imprime o resultado
		fmt.Println("cli: ", i, " req: ", v, " resp:", r)
	}
}

// ------------------------------------
// servidor
// thread de servico

func trataReq(id int, in chan Request){ 

	for {
		// espera requisicao de um cliente
		req := <-in
           //fmt.Println("Thread ",id)
		// responde a requisicao para o cliente
		req.ch_ret <- req.v * 2

		
	}
}


// ------------------------------------
// main
func main() {
	fmt.Println("Numero fixo de threads pre criadas, que leem o trabalho do canal de entrada")

	// cria canal entre clientes e servidor
	serv_chan := make(chan Request)
	
	// cria canal de sincronizacao final
	fim_chan := make(chan int)

	// cria numero fixo de threads
	for i := 0; i < MAXTH; i++{
		go trataReq(i, serv_chan)
	}
	
	// inicia a rotina dos clientes
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}
	<-fim_chan
}