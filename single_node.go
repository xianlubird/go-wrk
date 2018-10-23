package main

import (
	"sync"
	"time"
)

func SingleNode(toCall string) []byte {
	responseChannel := make(chan *Response, *totalCalls*2)

	benchTime := NewTimer()
	benchTime.Reset()
	//TODO check ulimit
	wg := &sync.WaitGroup{}

	for i := 0; i < *numConnections; i++ {
		go StartClient(
			toCall,
			*headers,
			*requestBody,
			*method,
			*disableKeepAlives,
			responseChannel,
			wg,
			*totalCalls,
		)
		wg.Add(1)
	}

	c := make(chan struct{})
	go func() {
        defer close(c)
        wg.Wait()
    }()

	if *timeNum <= 0 {
		*timeNum = 1
	}
	timeout := time.Duration(*timeNum) * time.Second

	select {
    case <-c:
    case <-time.After(timeout):
    }

	result := CalcStats(
		responseChannel,
		benchTime.Duration(),
	)
	return result
}
