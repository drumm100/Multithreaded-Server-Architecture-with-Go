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

// estrutura que define uma thread
type Thread struct {
	// identifica se a thread esta em uso
	uso bool	
	// define o canal entre servidor e a thread	
	thr chan Request 
}

var threads [MAXTH]Thread;
var indiceThreads int;

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

func achaCanal() Thread{

	// espera ate achar uma thread disponivel
	for {
		for i := 0; i < MAXTH; i++ {
			if (threads[i].uso == false) {
				threads[i].uso = true
				return threads[i]
			}
		}
	}

}

func trataReq(id int){ 
	var meuId int = id
	for {
		// espera requisicao enviada pelo servidor
		req := <- threads[meuId].thr
 
		// responde a requisicao para o cliente
		req.ch_ret <- req.v * 2

		// sinaliza que a thread esta livre
		threads[meuId].uso = false
	}
}

// servidor que dispara threads de servico
func servidor(in chan Request, fim chan int) {

	// cria numero fixo de threads
	for i := 0; i < MAXTH; i++{
		threads[i].uso = false
		threads[i].thr = make (chan Request)
		go trataReq(i)
	}
	
	for {
		// le requisicao do cliente
		req := <-in

		// acha uma thread disponivel para 
		// enviar a requisicao
		var canal Thread =  achaCanal()

		// envia requisicao para a thread
		canal.thr <- req
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("Numero fixo de threads pre criadas, que recebem trabalho de uma principal")

	// cria canal entre clientes e servidor
	serv_chan := make(chan Request)
	
	// cria canal de sincronizacao final
	fim_chan := make(chan int)

	// inicia a rotina do servidor
	go servidor(serv_chan, fim_chan)

	// inicia a rotina dos clientes
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}
	<-fim_chan
}