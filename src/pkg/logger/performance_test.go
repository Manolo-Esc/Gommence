package logger

import (
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// worker es la función de la goroutine, opcionalmente llama a logFunction
func worker(id int, useLog bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		val := math.Cos(float64(i))
		if useLog {
			GetLogger().Info(fmt.Sprintf("Goroutine %d ejecutando log: %f", id, val))
			//logFunction(id)
		}
	}
}

func benchmark(n int, useLog bool) time.Duration {
	var wg sync.WaitGroup
	wg.Add(n)

	start := time.Now()
	for i := 0; i < n; i++ {
		go worker(i, useLog, &wg)
	}
	wg.Wait()
	return time.Since(start)
}

func TestLoggingOverhead(t *testing.T) {
	_ = GetLogger() // create the singleton
	nThreads := 1000

	timeWithoutLog := benchmark(nThreads, false)
	timeWithLog := benchmark(nThreads, true)

	t.Logf("Tiempo sin logs: %v", timeWithoutLog)
	t.Logf("Tiempo con logs: %v", timeWithLog)
	t.Logf("Diferencia: %v", timeWithLog-timeWithoutLog)

	if timeWithLog > timeWithoutLog*2 {
		t.Errorf("El logging duplica el tiempo de ejecución")
	}
}

const bufferSize = 20000

type PriorityQueue struct {
	mu        sync.Mutex
	cond      *sync.Cond
	buffer    []*string
	maxSize   int
	available bool // Indica si hay datos para consumir
}

func NewPriorityQueue(size int) *PriorityQueue {
	pq := &PriorityQueue{
		buffer:  make([]*string, 0, size),
		maxSize: size,
	}
	pq.cond = sync.NewCond(&pq.mu)
	return pq
}

func (pq *PriorityQueue) Produce(str *string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	// Esperar solo si el buffer está lleno
	for len(pq.buffer) >= pq.maxSize {
		pq.cond.Wait() // Productores esperan solo si no hay espacio
	}

	pq.buffer = append(pq.buffer, str)
	pq.available = true
	pq.cond.Signal() // Notificar al consumidor si estaba esperando
}

func (pq *PriorityQueue) Consume(flushing bool) *string {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	// Esperar hasta que haya datos disponibles
	for !pq.available {
		if flushing {
			return nil
		}
		pq.cond.Wait()
	}

	item := pq.buffer[0]
	pq.buffer = pq.buffer[1:]

	// Si el buffer está vacío, notificar que no hay datos
	if len(pq.buffer) == 0 {
		pq.available = false
	}

	// Notificar a los productores que pueden escribir
	pq.cond.Broadcast()

	return item
}

func worker2(id int, useLog bool, wg *sync.WaitGroup, pq *PriorityQueue) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		val := math.Cos(float64(i))
		if useLog {
			msg := fmt.Sprintf("Goroutine %d ejecutando log: %f", id, val)
			pq.Produce(&msg)
			//logFunction(id)
		}
	}
}

func benchmark2(n int, useLog bool, pq *PriorityQueue) time.Duration {
	var wg sync.WaitGroup
	wg.Add(n)

	start := time.Now()
	for i := 0; i < n; i++ {
		go worker2(i, useLog, &wg, pq)
	}
	wg.Wait()
	return time.Since(start)
}

// func producer(id int, pq *PriorityQueue) {
// 	for i := 0; i < 5; i++ {
// 		item := id*10 + i
// 		fmt.Printf("Producer %d: Producing %d\n", id, item)
// 		pq.Produce(item)
// 		time.Sleep(time.Millisecond * 100) // Simula carga de trabajo
// 	}
// }

func consumer(pq *PriorityQueue, flag *atomic.Bool, file *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	flushing := false
	for {
		item := pq.Consume(flushing)
		if item == nil {
			return
		}
		if _, err := file.WriteString(*item + "\n"); err != nil {
			fmt.Println("Error escribiendo en el archivo:", err)
		}
		//fmt.Printf("Consumer: Consumed %d\n", item)
		//time.Sleep(time.Millisecond * 300) // Simula un consumidor más lento
		if !flushing && flag.Load() {
			flushing = true
		}
	}
}

func TestCustomQueue(t *testing.T) {
	var wgConsumer sync.WaitGroup
	wgConsumer.Add(1)

	file, err := os.OpenFile("miarchivo.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // os.O_TRUNC|os.O_CREATE|os.O_WRONLY
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)
		return
	}
	defer file.Close()

	var flag atomic.Bool
	pq := NewPriorityQueue(bufferSize)
	go consumer(pq, &flag, file, &wgConsumer)

	nThreads := 1000

	timeWithoutLog := benchmark2(nThreads, false, pq)
	timeWithLog := benchmark2(nThreads, true, pq)
	flag.Store(true)
	wgConsumer.Wait()

	t.Logf("Tiempo2 sin logs: %v", timeWithoutLog)
	t.Logf("Tiempo2 con logs: %v", timeWithLog)
	t.Logf("Diferencia2: %v", timeWithLog-timeWithoutLog)

	if timeWithLog > timeWithoutLog*2 {
		t.Errorf("El logging duplica el tiempo de ejecución")
	}
}
