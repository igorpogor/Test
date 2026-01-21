package main

import (
	feature1 "Study/feature"
	"context"
	"fmt"
	"time"
)

type Message struct {
	Author string
	Text   string
}

func DRUG1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("DRUG1: Вышел из сети")
			return
		default:
			fmt.Println("DRUG1: в сети")

		}
		time.Sleep(200 * time.Millisecond)

	}
}
func DRUG2(ctx context.Context) {

	for {

		select {
		case <-ctx.Done():
			fmt.Println("DRUG2: Вышел из сети")
			return
		default:
			fmt.Println("DRUG2: в сети")

		}
		time.Sleep(200 * time.Millisecond)

	}
}

func main() {
	parentContext, parentCancel := context.WithCancel(context.Background())
	childContext, childCancel := context.WithCancel(parentContext)
	go DRUG1(parentContext)
	go DRUG2(childContext)

	time.Sleep(2 * time.Second)
	childCancel()
	time.Sleep(2 * time.Second)
	parentCancel()
	time.Sleep(2 * time.Second)

	feature1.Feature1()
}
