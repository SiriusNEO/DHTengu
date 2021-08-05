package main

import (
	"flag"
	"kademlia"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	help     bool
	testName string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	kademlia.LogInit()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	_, _ = yellow.Println("Welcome to DHT-2020 Test Program!\n")

	var basicFailRate float64
	var forceQuitFailRate float64
	var QASFailRate float64

	_, _ = yellow.Println("Basic Test Begins:")
	basicPanicked, basicFailedCnt, basicTotalCnt := basicTest()
	if basicPanicked {
		_, _ = red.Printf("Basic Test Panicked.")
		os.Exit(0)
	}

	basicFailRate = float64(basicFailedCnt) / float64(basicTotalCnt)
	if basicFailRate > basicTestMaxFailRate {
		_, _ = red.Printf("Basic test failed with fail rate %.4f\n", basicFailRate)
	} else {
		_, _ = green.Printf("Basic test passed with fail rate %.4f\n", basicFailRate)
	}

	time.Sleep(afterTestSleepTime)

	_, _ = cyan.Println("\nFinal print:")
	if basicFailRate > basicTestMaxFailRate {
		_, _ = red.Printf("Basic test failed with fail rate %.4f\n", basicFailRate)
	} else {
		_, _ = green.Printf("Basic test passed with fail rate %.4f\n", basicFailRate)
	}
	if forceQuitFailRate > forceQuitMaxFailRate {
		_, _ = red.Printf("Force quit test failed with fail rate %.4f\n", forceQuitFailRate)
	} else {
		_, _ = green.Printf("Force quit test passed with fail rate %.4f\n", forceQuitFailRate)
	}
	if QASFailRate > QASMaxFailRate {
		_, _ = red.Printf("Quit & Stabilize test failed with fail rate %.4f\n", QASFailRate)
	} else {
		_, _ = green.Printf("Quit & Stabilize test passed with fail rate %.4f\n", QASFailRate)
	}
}

func usage() {
	flag.PrintDefaults()
}