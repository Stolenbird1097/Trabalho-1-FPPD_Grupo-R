// V1 – Versão concorrente sujeita a deadlock (solução "ingênua").
//
// Cada filósofo é uma goroutine e cada garfo é um sync.Mutex. Todos seguem a
// MESMA regra: pegam primeiro o garfo da esquerda e depois o da direita.
//
// Se os 5 pegarem o garfo esquerdo ao mesmo tempo, cada um fica segurando um
// garfo e esperando o da direita, que está com o vizinho:
//
//	F1 tem G0 e espera G1 → F2 tem G1 e espera G2 → ... → F5 tem G4 e espera G0
//
// Isso é uma espera circular e ninguém mais avança (deadlock). As quatro
// condições de Coffman estão presentes:
//   - exclusão mútua: um Mutex só pode estar com uma goroutine;
//   - segura e espera: o filósofo segura o garfo esquerdo enquanto espera o direito;
//   - sem preempção: ninguém tira um garfo de outro; só o dono faz Unlock;
//   - espera circular: a cadeia F1 → F2 → ... → F5 → F1 acima.
//
// Para a demonstração ser reproduzível, há uma pequena pausa entre pegar o
// primeiro e o segundo garfo (flag -janela, padrão 1 ms). Ela não muda a
// lógica, só aumenta a chance de todos estarem com o garfo esquerdo ao mesmo
// tempo. Com -janela 0 o deadlock continua possível, mas é mais raro.
//
// Um "vigia" (watchdog) percebe quando ninguém avança por 2 s, imprime quem
// segura qual garfo e quem espera qual (a evidência pedida no enunciado) e
// encerra o programa com código de saída 2.
//
// Execução:  go run ./cmd/v1-deadlock  [-v] [-csv] [-janela 1]
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SEU-USUARIO/filosofos-jantando/internal/config"
	"github.com/SEU-USUARIO/filosofos-jantando/internal/metrics"
)

const semGarfo = -1

// estado guarda, para a evidência do deadlock, o que cada filósofo segura e
// por qual garfo está esperando. É protegido pelo próprio mutex.
type estado struct {
	mu        sync.Mutex
	segura    [config.NumFilosofos][]int
	esperando [config.NumFilosofos]int
}

func (e *estado) quer(id, garfo int) {
	e.mu.Lock()
	e.esperando[id] = garfo
	e.mu.Unlock()
}

func (e *estado) pegou(id, garfo int) {
	e.mu.Lock()
	e.esperando[id] = semGarfo
	e.segura[id] = append(e.segura[id], garfo)
	e.mu.Unlock()
}

func (e *estado) soltouTudo(id int) {
	e.mu.Lock()
	e.segura[id] = nil
	e.mu.Unlock()
}

func (e *estado) relatorio() {
	e.mu.Lock()
	defer e.mu.Unlock()
	fmt.Println("\n*** DEADLOCK DETECTADO: nenhum filósofo avançou nos últimos 2 s ***")
	fmt.Println("Filósofo | Segura garfo | Espera garfo | Garfo esperado está com")
	for id := 0; id < config.NumFilosofos; id++ {
		dono := "-"
		if g := e.esperando[id]; g != semGarfo {
			for outro := 0; outro < config.NumFilosofos; outro++ {
				for _, s := range e.segura[outro] {
					if s == g {
						dono = fmt.Sprintf("Filósofo %d", outro+1)
					}
				}
			}
		}
		segura := "-"
		for i, g := range e.segura[id] {
			if i == 0 {
				segura = garfoTexto(g)
			} else {
				segura += ", " + garfoTexto(g)
			}
		}
		fmt.Printf("   F%d    | %-12s | %-12s | %s\n",
			id+1, segura, garfoTexto(e.esperando[id]), dono)
	}
	fmt.Println("Cada filósofo segura um garfo e espera o que está com o vizinho:")
	fmt.Println("espera circular F1 → F2 → F3 → F4 → F5 → F1.")
}

func garfoTexto(g int) string {
	if g == semGarfo {
		return "-"
	}
	return fmt.Sprintf("G%d", g)
}

func main() {
	janela := flag.Int("janela", 1, "pausa em ms entre pegar o 1º e o 2º garfo")
	cfg := config.Carregar() // chama flag.Parse(), que também lê -janela
	m := metrics.Nova("V1", config.NumFilosofos)

	var garfos [config.NumFilosofos]sync.Mutex
	st := &estado{}
	for i := range st.esperando {
		st.esperando[i] = semGarfo
	}
	var progresso atomic.Int64 // incrementado a cada ciclo concluído

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

			// 2. Adquirir: SEMPRE esquerdo primeiro, depois direito
			inicio := time.Now()
			st.quer(id, esq)
			garfos[esq].Lock()
			st.pegou(id, esq)
			log("F%d pegou G%d (esquerdo) e agora quer G%d", id+1, esq, dir)

			time.Sleep(time.Duration(*janela) * time.Millisecond)

			st.quer(id, dir)
			garfos[dir].Lock() // <- aqui todos podem ficar presos
			st.pegou(id, dir)
			m.RegistrarEspera(id, time.Since(inicio))

			// 3. Comer
			log("F%d comendo com G%d e G%d", id+1, esq, dir)
			m.ComecouAComer()
			time.Sleep(cfg.TempoComer)
			m.TerminouDeComer()

			// 4. Liberar
			garfos[dir].Unlock()
			garfos[esq].Unlock()
			st.soltouTudo(id)
			m.CicloConcluido(id)
			progresso.Add(1)
		}
	}

	m.Iniciar()
	var wg sync.WaitGroup
	for id := 0; id < config.NumFilosofos; id++ {
		wg.Add(1)
		go filosofo(id, &wg)
	}

	// Espera o fim normal ou detecta o deadlock.
	terminou := make(chan struct{})
	go func() { wg.Wait(); close(terminou) }()

	ultimo := int64(-1)
	vigia := time.NewTicker(2 * time.Second)
	defer vigia.Stop()
	for {
		select {
		case <-terminou:
			m.Finalizar()
			if cfg.CSV {
				fmt.Println(m.LinhaCSV())
			} else {
				fmt.Println("Execução terminou SEM deadlock desta vez.")
				m.Imprimir()
			}
			return
		case <-vigia.C:
			atual := progresso.Load()
			if atual == ultimo {
				m.Finalizar()
				st.relatorio()
				if !cfg.CSV {
					m.Imprimir() // mostra quantos ciclos cada um completou antes de travar
				}
				os.Exit(2)
			}
			ultimo = atual
		}
	}
}
