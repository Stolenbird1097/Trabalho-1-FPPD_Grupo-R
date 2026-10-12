// V2 – Sem deadlock, mas com sincronização excessiva ("garçom único").
//
// Cada filósofo é uma goroutine e cada garfo é um sync.Mutex, como na V1.
// A diferença é um SEMÁFORO BINÁRIO global (o "garçom"): para pegar os
// garfos o filósofo precisa antes da permissão do garçom, e só a devolve
// DEPOIS de comer e soltar os garfos. Ou seja, a região crítica é a refeição
// inteira, e só um filósofo come por vez em toda a mesa.
//
// Por que não há deadlock: como só um filósofo pode estar adquirindo/usando
// garfos por vez, ele nunca encontra um garfo ocupado. Isso quebra a
// condição de "segura e espera" (quem tem a permissão não espera por
// ninguém) e, portanto, a espera circular nunca se forma.
//
// Por que é excessiva: os filósofos F1 e F3, por exemplo, não dividem
// nenhum garfo e poderiam comer juntos, mas o garçom não deixa. O máximo de
// filósofos comendo simultaneamente é sempre 1, e quem chega precisa esperar
// a refeição de todos que estão na fila.
//
// Mecanismo sem equivalente nativo em Go: Go não tem um tipo "semáforo" na
// biblioteca padrão. Ele foi adaptado com um canal com buffer de capacidade
// 1 (chan struct{}): enviar no canal = P()/wait (bloqueia se cheio);
// receber do canal = V()/signal (libera a vaga). A capacidade do buffer é o
// valor inicial do semáforo.
//
// Execução:  go run ./cmd/v2-sincronizacao-excessiva  [-v] [-csv]
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/SEU-USUARIO/filosofos-jantando/internal/config"
	"github.com/SEU-USUARIO/filosofos-jantando/internal/metrics"
)

// Semaforo é um semáforo contador implementado com canal com buffer.
type Semaforo chan struct{}

func NovoSemaforo(vagas int) Semaforo { return make(Semaforo, vagas) }
func (s Semaforo) P()                 { s <- struct{}{} } // wait / down
func (s Semaforo) V()                 { <-s }             // signal / up

func main() {
	cfg := config.Carregar()
	m := metrics.Nova("V2", config.NumFilosofos)

	var garfos [config.NumFilosofos]sync.Mutex
	garcom := NovoSemaforo(1) // semáforo binário: 1 filósofo por vez

	log := func(formato string, args ...any) {
		if cfg.Verbose {
			fmt.Printf(formato+"\n", args...)
		}
	}

	filosofo := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()
		esq, dir := id, (id+1)%config.NumFilosofos

		for ciclo := 1; ciclo <= config.NumCiclos; ciclo++ {
			// 1. Pensar
			log("F%d pensando (ciclo %d)", id+1, ciclo)
			time.Sleep(cfg.TempoPensar)

			// 2. Adquirir: permissão do garçom + os dois garfos
			inicio := time.Now()
			garcom.P() // espera aqui enquanto outro filósofo come
			garfos[esq].Lock()
			garfos[dir].Lock()
			m.RegistrarEspera(id, time.Since(inicio))

			// 3. Comer (ainda segurando a permissão do garçom)
			log("F%d comendo com G%d e G%d", id+1, esq, dir)
			m.ComecouAComer()
			time.Sleep(cfg.TempoComer)
			m.TerminouDeComer()

			// 4. Liberar garfos e, só então, a permissão
			garfos[dir].Unlock()
			garfos[esq].Unlock()
			garcom.V()
			m.CicloConcluido(id)
		}
	}

	m.Iniciar()
	var wg sync.WaitGroup
	for id := 0; id < config.NumFilosofos; id++ {
		wg.Add(1)
		go filosofo(id, &wg)
	}
	wg.Wait()
	m.Finalizar()

	if cfg.CSV {
		fmt.Println(m.LinhaCSV())
	} else {
		m.Imprimir()
	}
}
