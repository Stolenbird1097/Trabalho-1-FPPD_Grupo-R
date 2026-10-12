// V3 - Concurrent Dining Philosophers without deadlock (resource ordering)
// Same as V1, with one change: every philosopher takes the fork with the
// LOWER number first. Philosophers 1-4 behave like V1, but philosopher 5
// (forks 4 and 0) takes fork 0 first. This breaks the circular wait,
// so there is no deadlock, and up to 2 philosophers can eat together.
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
	first := id
	second := (id + 1) % philosophers
	if first > second { // always the lower fork number first
		first, second = second, first
	}

	for i := 1; i <= rounds; i++ {
		// think
		show(id, fmt.Sprintf(" T%d", i))
		time.Sleep(time_to_think * time.Millisecond)

		// get forks: lower number first
		start := time.Now()
		forks[first].Lock()
		forks[second].Lock()
		startEating(time.Since(start))

		// eat
		show(id, fmt.Sprintf(" E%d", i))
		time.Sleep(time_to_eat * time.Millisecond)

		// release forks
		forks[second].Unlock()
		forks[first].Unlock()
		stopEating(id)
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
