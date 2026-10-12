// V1 - Jantar dos Filósofos concorrente COM deadlock
// Cada filósofo é uma goroutine e cada garfo é um sync.Mutex.
// Todos pegam primeiro o garfo da ESQUERDA e depois o da DIREITA.
// Se os 5 pegarem o garfo esquerdo ao mesmo tempo, todos esperam para sempre
// pelo garfo direito (espera circular) -> deadlock.
// O Go detecta e encerra com "fatal error: all goroutines are asleep - deadlock!"
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
	left := id
	right := (id + 1) % philosophers

	for i := 1; i <= rounds; i++ {
		// pensar
		show(id, fmt.Sprintf(" T%d", i))
		time.Sleep(time_to_think * time.Millisecond)

		// pegar os garfos: primeiro o esquerdo, depois o direito
		start := time.Now()
		forks[left].Lock()
		show(id, fmt.Sprintf(" G%d", left)) // segurando o garfo esquerdo
		time.Sleep(1 * time.Millisecond)    // leva um instante para alcançar o outro garfo
		forks[right].Lock()                 // <- aqui todos podem ficar presos
		startEating(time.Since(start))

		// comer
		show(id, fmt.Sprintf(" E%d", i))
		time.Sleep(time_to_eat * time.Millisecond)
		stopEating(id) // terminou de comer (antes de soltar os garfos, para a contagem ser exata)

		// soltar os garfos
		forks[right].Unlock()
		forks[left].Unlock()
	}
}

func main() {
	fmt.Println("\n [P1]  [P2]  [P3]  [P4]  [P5]")
	fmt.Println(" T = thinking, G = holding fork, E = eating")

	var wg sync.WaitGroup
	start := time.Now()
	for j := 0; j < philosophers; j++ {
		wg.Add(1)
		go Philosopher(j, &wg)
	}
	wg.Wait()
	report("V1", time.Since(start))
}
