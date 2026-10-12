// V3 – Sem deadlock e com mais concorrência (ordenação de recursos).
//
// Cada filósofo é uma goroutine e cada garfo é um sync.Mutex, como na V1.
// A única mudança é a ORDEM de aquisição: todo filósofo pega primeiro o
// garfo de MENOR número e depois o de maior número (hierarquia de recursos,
// proposta por Dijkstra).
//
// Para F1..F4 isso é igual à V1 (esquerdo, depois direito). Já F5, cujos
// garfos são G4 (esquerdo) e G0 (direito), passa a pegar G0 primeiro. Assim
// F1 e F5 disputam G0 como PRIMEIRO garfo: quem perde fica esperando sem
// segurar nada.
//
// Por que não há deadlock: numa espera circular, cada filósofo esperaria um
// garfo de número maior do que o que já segura (G_a < G_b < ... < G_a), o que
// é impossível. A condição de espera circular é eliminada, e sem ela não há
// deadlock, mesmo mantendo exclusão mútua, segura-e-espera e não preempção.
//
// Concorrência: não há trava global; quem não divide garfos come ao mesmo
// tempo. Com 5 filósofos e 5 garfos, até 2 comem juntos (5/2 arredondado
// para baixo), que é o máximo possível no problema.
//
// Starvation: em teoria é possível (nada garante que um filósofo vença a
// disputa de um garfo). Na prática é muito improvável aqui, porque o
// sync.Mutex do Go entra em "modo de inanição" quando uma goroutine espera
// mais de 1 ms e passa a entregar o lock em ordem FIFO, e porque cada
// filósofo pensa entre as refeições, dando chance aos vizinhos. A métrica
// "ciclos por filósofo" (todos com 20) e o "maior tempo de espera" mostram
// isso experimentalmente.
//
// Execução:  go run ./cmd/v3-concorrente  [-v] [-csv]
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/SEU-USUARIO/filosofos-jantando/internal/config"
	"github.com/SEU-USUARIO/filosofos-jantando/internal/metrics"
)

func main() {
	cfg := config.Carregar()
	m := metrics.Nova("V3", config.NumFilosofos)

	var garfos [config.NumFilosofos]sync.Mutex

	log := func(formato string, args ...any) {
		if cfg.Verbose {
			fmt.Printf(formato+"\n", args...)
		}
	}

	filosofo := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()
		esq, dir := id, (id+1)%config.NumFilosofos

		// Ordem global: menor número primeiro.
		primeiro, segundo := esq, dir
		if primeiro > segundo {
			primeiro, segundo = segundo, primeiro
		}

		for ciclo := 1; ciclo <= config.NumCiclos; ciclo++ {
			// 1. Pensar
			log("F%d pensando (ciclo %d)", id+1, ciclo)
			time.Sleep(cfg.TempoPensar)

			// 2. Adquirir na ordem global
			inicio := time.Now()
			garfos[primeiro].Lock()
			garfos[segundo].Lock()
			m.RegistrarEspera(id, time.Since(inicio))

			// 3. Comer
			log("F%d comendo com G%d e G%d", id+1, primeiro, segundo)
			m.ComecouAComer()
			time.Sleep(cfg.TempoComer)
			m.TerminouDeComer()

			// 4. Liberar (ordem inversa, por convenção)
			garfos[segundo].Unlock()
			garfos[primeiro].Unlock()
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
