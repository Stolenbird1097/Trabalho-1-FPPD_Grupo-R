// Package metrics coleta as métricas exigidas pelo enunciado, de forma igual
// para todas as versões (V0 a V3):
//
//   - tempo total de execução
//   - tempo médio de espera para aquisição dos garfos
//   - maior tempo de espera observado
//   - número máximo de filósofos comendo simultaneamente
//   - número de ciclos concluídos por cada filósofo
//
// Uso típico dentro do ciclo de um filósofo:
//
//	inicio := time.Now()
//	// ... adquirir os dois garfos ...
//	m.RegistrarEspera(id, time.Since(inicio))
//	m.ComecouAComer()
//	time.Sleep(cfg.TempoComer)
//	m.TerminouDeComer()
//	// ... liberar os garfos ...
//	m.CicloConcluido(id)
package metrics

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Metricas struct {
	versao string
	inicio time.Time
	fim    time.Time

	mu          sync.Mutex
	somaEspera  time.Duration
	qtdEsperas  int
	maiorEspera time.Duration
	ciclos      []int

	comendoAgora  atomic.Int32
	maxSimultaneo atomic.Int32
}

func Nova(versao string, numFilosofos int) *Metricas {
	return &Metricas{versao: versao, ciclos: make([]int, numFilosofos)}
}

func (m *Metricas) Iniciar()   { m.inicio = time.Now() }
func (m *Metricas) Finalizar() { m.fim = time.Now() }

// RegistrarEspera deve receber o tempo entre "quero comer" e "tenho os dois garfos".
func (m *Metricas) RegistrarEspera(id int, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.somaEspera += d
	m.qtdEsperas++
	if d > m.maiorEspera {
		m.maiorEspera = d
	}
}

func (m *Metricas) ComecouAComer() {
	agora := m.comendoAgora.Add(1)
	for {
		max := m.maxSimultaneo.Load()
		if agora <= max || m.maxSimultaneo.CompareAndSwap(max, agora) {
			return
		}
	}
}

func (m *Metricas) TerminouDeComer() { m.comendoAgora.Add(-1) }

func (m *Metricas) CicloConcluido(id int) {
	m.mu.Lock()
	m.ciclos[id]++
	m.mu.Unlock()
}

func (m *Metricas) Total() time.Duration { return m.fim.Sub(m.inicio) }

func (m *Metricas) EsperaMedia() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.qtdEsperas == 0 {
		return 0
	}
	return m.somaEspera / time.Duration(m.qtdEsperas)
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 }

// CabecalhoCSV é a primeira linha do arquivo gerado por scripts/executar.sh.
const CabecalhoCSV = "versao,total_ms,espera_media_ms,espera_max_ms,max_simultaneos,ciclos_por_filosofo"

func (m *Metricas) LinhaCSV() string {
	m.mu.Lock()
	ciclos := make([]string, len(m.ciclos))
	for i, c := range m.ciclos {
		ciclos[i] = fmt.Sprint(c)
	}
	maior := m.maiorEspera
	m.mu.Unlock()
	return fmt.Sprintf("%s,%.3f,%.3f,%.3f,%d,%s",
		m.versao, ms(m.Total()), ms(m.EsperaMedia()), ms(maior),
		m.maxSimultaneo.Load(), strings.Join(ciclos, "|"))
}

func (m *Metricas) Imprimir() {
	m.mu.Lock()
	maior := m.maiorEspera
	ciclos := append([]int(nil), m.ciclos...)
	m.mu.Unlock()
	fmt.Printf("\n=== Resultados %s ===\n", m.versao)
	fmt.Printf("Tempo total de execução ........ %.3f ms\n", ms(m.Total()))
	fmt.Printf("Tempo médio de espera (garfos) . %.3f ms\n", ms(m.EsperaMedia()))
	fmt.Printf("Maior tempo de espera .......... %.3f ms\n", ms(maior))
	fmt.Printf("Máx. comendo simultaneamente ... %d\n", m.maxSimultaneo.Load())
	for i, c := range ciclos {
		fmt.Printf("Filósofo %d: %d ciclos\n", i+1, c)
	}
}
