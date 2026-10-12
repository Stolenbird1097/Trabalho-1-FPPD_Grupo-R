// V2 - Dining Philosophers without deadlock, but with EXCESSIVE synchronization
// Same as V1, plus a "waiter": a semaphore with only 1 permit.
// A philosopher asks the waiter before taking the forks and only gives
// the permit back after eating. So only ONE philosopher eats at a time,
// even when two of them don't share any fork.
// Go has no native semaphore: it is made with a buffered channel of size 1
// (send = wait/P, receive = signal/V).
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	philosophers  = 5
	rounds        = 20
	time_to_think = 20 // ms (same in V0, V1, V2, V3)
	time_to_eat   = 30 // ms (same in V0, V1, V2, V3)
)

var forks [philosophers]sync.Mutex
var waiter = make(chan struct{}, 1) // semaphore with 1 permit

// ---------- metrics (same code in all versions) ----------
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
	// last line in CSV (used by scripts/executar.sh)
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
		// think
		show(id, fmt.Sprintf(" T%d", i))
		time.Sleep(time_to_think * time.Millisecond)

		// get forks: first the waiter's permit, then left and right
		start := time.Now()
		waiter <- struct{}{} // wait (P): blocks while someone else is eating
		forks[left].Lock()
		forks[right].Lock()
		startEating(time.Since(start))

		// eat
		show(id, fmt.Sprintf(" E%d", i))
		time.Sleep(time_to_eat * time.Millisecond)

		// release forks and only then the permit
		forks[right].Unlock()
		forks[left].Unlock()
		stopEating(id)
		<-waiter // signal (V)
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
	report("V2", time.Since(start))
}
