package main

import (
	"fmt"
)

type (
	Bytes   int
	Celcius float32
)

type BandwidthUsage struct {
	amount []Bytes
}

type CpuTemp struct {
	temp []Celcius
}

type MemoryUsage struct {
	amount []Bytes
}

func (B *BandwidthUsage) getBandwidthAverage() Bytes {
	var average Bytes
	allItems := len(B.amount)
	for _, amount := range B.amount {
		average += amount
	}
	average = average / Bytes(allItems)
	return average
}

func (C *CpuTemp) getTempAverage() Celcius {
	var average Celcius
	allItems := len(C.temp)
	for _, temp := range C.temp {
		average += temp
	}
	average = average / Celcius(allItems)
	return average
}

func (M *MemoryUsage) getMemoryUsageAverage() Bytes {
	var average Bytes
	allItems := len(M.amount)
	for _, amount := range M.amount {
		average += amount
	}
	average = average / Bytes(allItems)
	return average
}

func Compare[T comparable](n, n1 T) bool {
	return n == n1
}

func main() {
	bandwidthUsage := BandwidthUsage{[]Bytes{5000, 300, 500, 200, 30}}
	cpuTemp := CpuTemp{[]Celcius{50, 40, 10, 40, 54}}
	memory := MemoryUsage{[]Bytes{40, 60, 30, 70, 20}}

	fmt.Println(bandwidthUsage.getBandwidthAverage())
	fmt.Println(cpuTemp.getTempAverage())
	fmt.Println(memory.getMemoryUsageAverage())
}
