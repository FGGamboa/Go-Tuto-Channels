package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Types of channels
// bidireccionales
// receive-only
// send-only

func main() {
	wg := &sync.WaitGroup{}
	IdsChan := make(chan string)
	FakeIdsChan := make(chan string)
	ClosedChans := make(chan int)

	wg.Add(3)

	go generateIds(wg, IdsChan, ClosedChans)
	go generateFakeIds(wg, FakeIdsChan, ClosedChans)
	go logIds(wg, IdsChan, FakeIdsChan, ClosedChans)

	wg.Wait()
}

func generateFakeIds(wg *sync.WaitGroup, fakeIdsChan chan<- string, closedChannels chan<- int) {
	for i := 0; i < 50; i++ {
		id := uuid.New()
		fakeIdsChan <- fmt.Sprintf("%d . %s", i+1, id.String())
	}
	close(fakeIdsChan)
	closedChannels <- 1

	wg.Done()
}

// Recibiendo datos <- derecha
func generateIds(wg *sync.WaitGroup, idsChan chan<- string, closedChannels chan<- int) {
	for i := 0; i < 100; i++ {
		id := uuid.New()
		idsChan <- fmt.Sprintf("%d . %s", i+1, id.String())
	}

	close(idsChan)
	closedChannels <- 1

	wg.Done()
}

// Escuchando datos <- izquierda
func logIds(wg *sync.WaitGroup, idsChan <-chan string, fakeIdsChan <-chan string, closedChannels chan int) {
	closedCounter := 0

	for {
		select {
		case id, ok := <-idsChan:
			if ok {
				fmt.Println("Id: ", id)
			}

		case id, ok := <-fakeIdsChan:
			if ok {
				fmt.Println("FakeId: ", id)
			}

		case count, ok := <-closedChannels:
			if ok {
				closedCounter += count
			}
		}

		if closedCounter == 2 {
			close(closedChannels)
			break
		}
	}

	wg.Done()
}
