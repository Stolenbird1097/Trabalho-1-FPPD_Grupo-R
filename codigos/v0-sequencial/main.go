// V0 - Jantar dos Filósofos sequencial (baseline)
// Baseado em filo_seq.go de Cesar De Rose (Moodle).
// Um filósofo por vez: sem concorrência e sem disputa pelos garfos.
// Tempo esperado = rounds * philosophers * (time_to_think + time_to_eat)
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	philosophers  = 5
	rounds        = 20
	time_to_think = 20 // ms (igual em V0, V1, V2, V3)
	time_to_eat   = 30 // ms (igual em V0, V1, V2, V3)
)

// ---------- métricas (mesmo código em todas as versões) ----------
var (
	mu                sync.Mutex
	totalWait         time.Duration
	maxWait           time.Duration
	waits             int
	eating, maxEating int
	cycles            [philosophers]int
)

func startEating(wait time.Duration) {
	mu.Lock()
	totalWait += wait
	waits++
	if wait > maxWait {
		maxWait = wait
	}
	eating++
	if eating > maxEating {
		maxEating = eating
	}
	mu.Unlock()
}

func stopEating(id int) {
	mu.Lock()
	eating--
	cycles[id]++
	mu.Unlock()
}

func report(version string, elapsed time.Duration) {
	avg := totalWait / time.Duration(waits)
	fmt.Printf("\n%s took %s\n", version, elapsed)
	fmt.Printf("Average wait for forks: %s\n", avg)
	fmt.Printf("Max wait for forks:     %s\n", maxWait)
	fmt.Printf("Max eating together:    %d\n", maxEating)
	fmt.Printf("Cycles per philosopher: %v\n", cycles)
	// última linha em CSV (usada por scripts/executar.sh)
	ms := func(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 }
	fmt.Printf("%s,%.3f,%.3f,%.3f,%d,%d|%d|%d|%d|%d\n", version, ms(elapsed), ms(avg), ms(maxWait),
		maxEating, cycles[0], cycles[1], cycles[2], cycles[3], cycles[4])
}

func show(id int, text string) { fmt.Printf("%*s%s\n", id*6, "", text) }

// ---------------------------------------------------------

func Philosopher(id, iteration int) {
	// pensar
	show(id, fmt.Sprintf(" T%d", iteration))
	time.Sleep(time_to_think * time.Millisecond)

	// pegar os garfos (ninguém mais está executando, então não há espera)
	start := time.Now()
	startEating(time.Since(start))

	// comer
	show(id, fmt.Sprintf(" E%d", iteration))
	time.Sleep(time_to_eat * time.Millisecond)

	// soltar os garfos
	stopEating(id)
}

func main() {
	fmt.Println("\n [P1]  [P2]  [P3]  [P4]  [P5]")

	start := time.Now()
	for i := 1; i <= rounds; i++ {
		for j := 0; j < philosophers; j++ {
			Philosopher(j, i)
		}
	}
	report("V0", time.Since(start))
}
