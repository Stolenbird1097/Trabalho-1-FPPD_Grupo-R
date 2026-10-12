// V0 – Versão sequencial (baseline).
//
// Os filósofos NÃO executam ao mesmo tempo: existe uma única linha de
// execução (a goroutine main) que faz, em rodízio, um ciclo completo de cada
// filósofo:
//
//	pensar → adquirir garfos → comer → liberar garfos
//
// Como só um filósofo está ativo por vez, nunca há disputa pelos garfos:
// não existe deadlock, starvation nem condição de corrida, e no máximo 1
// filósofo come por vez. O tempo total é a soma de todos os tempos de
// pensar e comer (5 filósofos × 20 ciclos × (pensar + comer)), por isso esta
// versão serve de referência para medir o ganho das versões concorrentes.
//
// Execução:  go run ./cmd/v0-sequencial  [-v] [-csv]
package main

import (
	"fmt"
	"time"

	"github.com/SEU-USUARIO/filosofos-jantando/internal/config"
	"github.com/SEU-USUARIO/filosofos-jantando/internal/metrics"
)

// Na versão sequencial um garfo é só um booleano "em uso": como não há
// concorrência, não é preciso nenhum mecanismo de sincronização.
var garfos [config.NumFilosofos]bool

func esquerdo(id int) int { return id }
func direito(id int) int  { return (id + 1) % config.NumFilosofos }

func main() {
	cfg := config.Carregar()
	m := metrics.Nova("V0", config.NumFilosofos)

	log := func(formato string, args ...any) {
		if cfg.Verbose {
			fmt.Printf(formato+"\n", args...)
		}
	}

	m.Iniciar()
	for ciclo := 1; ciclo <= config.NumCiclos; ciclo++ {
		for id := 0; id < config.NumFilosofos; id++ {
			// 1. Pensar
			log("[ciclo %2d] Filósofo %d pensando", ciclo, id+1)
			time.Sleep(cfg.TempoPensar)

			// 2. Adquirir os garfos (nunca espera: ninguém mais está ativo)
			inicio := time.Now()
			garfos[esquerdo(id)] = true
			garfos[direito(id)] = true
			m.RegistrarEspera(id, time.Since(inicio))

			// 3. Comer
			log("[ciclo %2d] Filósofo %d comendo com garfos %d e %d",
				ciclo, id+1, esquerdo(id), direito(id))
			m.ComecouAComer()
			time.Sleep(cfg.TempoComer)
			m.TerminouDeComer()

			// 4. Liberar os garfos
			garfos[esquerdo(id)] = false
			garfos[direito(id)] = false
			m.CicloConcluido(id)
		}
	}
	m.Finalizar()

	if cfg.CSV {
		fmt.Println(m.LinhaCSV())
	} else {
		m.Imprimir()
	}
}
