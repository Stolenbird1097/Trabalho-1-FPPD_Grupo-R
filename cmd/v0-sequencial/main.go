// V0 - Sequential Dining Philosophers (baseline)
// Based on filo_seq.go by Cesar De Rose (Moodle).
// One philosopher at a time: no concurrency, no forks dispute.
// Expected time = rounds * philosophers * (time_to_think + time_to_eat)
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

func Philosopher(id, iteration int) {
	// think
	show(id, fmt.Sprintf(" T%d", iteration))
	time.Sleep(time_to_think * time.Millisecond)

	// get forks (nobody else is running, so no waiting)
	start := time.Now()
	startEating(time.Since(start))

	// eat
	show(id, fmt.Sprintf(" E%d", iteration))
	time.Sleep(time_to_eat * time.Millisecond)

	// release forks
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
