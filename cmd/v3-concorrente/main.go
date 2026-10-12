// V3 - Jantar dos Filósofos concorrente sem deadlock (ordenação de recursos)
// Igual à V1, com uma mudança: todo filósofo pega primeiro o garfo de
// MENOR número. Os filósofos 1-4 se comportam como na V1, mas o filósofo 5
// (garfos 4 e 0) pega o garfo 0 primeiro. Isso quebra a espera circular,
// então não há deadlock, e até 2 filósofos podem comer juntos.
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

var forks [philosophers]sync.Mutex

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

func Philosopher(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	first := id
	second := (id + 1) % philosophers
	if first > second { // sempre o garfo de menor número primeiro
		first, second = second, first
	}

	for i := 1; i <= rounds; i++ {
		// pensar
		show(id, fmt.Sprintf(" T%d", i))
		time.Sleep(time_to_think * time.Millisecond)

		// pegar os garfos: menor número primeiro
		start := time.Now()
		forks[first].Lock()
		forks[second].Lock()
		startEating(time.Since(start))

		// comer
		show(id, fmt.Sprintf(" E%d", i))
		time.Sleep(time_to_eat * time.Millisecond)
		stopEating(id) // terminou de comer (antes de soltar os garfos, para a contagem ser exata)

		// soltar os garfos
		forks[second].Unlock()
		forks[first].Unlock()
	}
}

func main() {
	fmt.Println("\n [P1]  [P2]  [P3]  [P4]  [P5]")

	var wg sync.WaitGroup
	start := time.Now()
	for j := 0; j < philosophers; j++ {
		wg.Add(1)
		go Philosopher(j, &wg)
	}
	wg.Wait()
	report("V3", time.Since(start))
}
