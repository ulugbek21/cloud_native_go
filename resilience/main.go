package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"

	circuitbreaker "github.com/ulugbek21/cloud_native_go/resilience/circuit_breaker"
	"github.com/ulugbek21/cloud_native_go/resilience/debounce"
	"github.com/ulugbek21/cloud_native_go/resilience/retry"
)

var fl = flag.String("t", "circuit", "type of pattern - default \"circuit\"")

func main() {
	flag.Parse()

	if *fl == "circuit" {
		// Circuit breaker
		fmt.Println("Circuit breaker output:")
		var circuit circuitbreaker.Circuit = func(c context.Context) (string, error) {
			return "", errors.New("error")
		}

		circuit = circuitbreaker.Breaker(circuit, 2)

		fmt.Println(circuit(context.Background()))
		fmt.Println(circuit(context.Background()))
		fmt.Println(circuit(context.Background()))
		<-time.After(time.Second * 3)
		fmt.Println(circuit(context.Background()))
		fmt.Println()
	}

	if *fl == "debounce" {
		// Debounce
		fmt.Println("Debounce output:")
		var debounceCircuit debounce.Circuit = func(c context.Context) (string, error) {
			return strconv.Itoa(int(time.Now().Unix())), nil
		}

		// Debounce first (debounce last works same but for last invocation)
		debounceCircuit = debounce.DebounceFirst(debounceCircuit, time.Second*1)
		fmt.Println(debounceCircuit(context.Background()))
		<-time.After(time.Millisecond * 500)
		fmt.Println("After 500 milliseconds. No change in value expected:")
		fmt.Println(debounceCircuit(context.Background()))
		<-time.After(time.Millisecond * 1200)
		fmt.Println("After 1.7 seconds. Change in value is expected:")
		fmt.Println(debounceCircuit(context.Background()))
	}

	if *fl == "retry" {
		// Retry
		effector := func(c context.Context) (string, error) {
			return "", errors.New("retry error")
		}

		effector = retry.Retry(effector, 2, time.Millisecond*100)

		fmt.Println(effector(context.Background()))
	}
}
