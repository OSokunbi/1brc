package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Station struct {
	min     float64
	max     float64
	average float64
	visited int
}

func main() {
	data := make(map[string]Station)
	file, err := os.Open("measurements.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	var mutex sync.Mutex
	start := time.Now()
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			res := strings.Split(line, ";")
			station := res[0]
			weather, err := strconv.ParseFloat(res[1], 64)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			mutex.Lock()
			defer mutex.Unlock()

			if _, ok := data[station]; !ok {
				data[station] = Station{
					min:     weather,
					max:     weather,
					average: weather,
					visited: 1,
				}
			} else {
				data[station] = Station{
					min:     min(data[station].min, weather),
					max:     max(data[station].max, weather),
					average: (data[station].average*float64(data[station].visited) + weather) / float64(data[station].visited+1),
					visited: data[station].visited + 1,
				}
			}
		}(line)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	
	elapsed := time.Since(start)
	fmt.Printf("Elapsed time %s\n", elapsed)
	fmt.Println("Chicago:", data["Chicago"])
}

