package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Measurement struct {
		minTemperature float64
	 	maxTemperature float64
	  sumTemperature float64
		countTemperature int64
} 

func main() {
	  start := time.Now()
		const k = 100_000_000
		var wg sync.WaitGroup
		wg.Add(k)
		 
	  measurements, err := os.Open("measurements.txt")
		if err != nil {
			  panic(err)
		}
		defer measurements.Close()

		data := make(map[string]Measurement)
		
		scanner := bufio.NewScanner(measurements)
		for scanner.Scan() {
			rawData := scanner.Text()
			semicolon := strings.Index(rawData, ";") 
			location := rawData[:semicolon]
			rawTemp := rawData[semicolon+1:]

			temperature, _ := strconv.ParseFloat(rawTemp, 64)

			measurement, ok := data[location]
			if !ok {
				  measurement = Measurement{
							minTemperature: temperature,
							maxTemperature: temperature,
							sumTemperature: temperature,
							countTemperature: 1,
						}
      } else {
				  measurement.minTemperature = min(measurement.minTemperature, temperature)
			    measurement.maxTemperature = max(measurement.maxTemperature, temperature)
			    measurement.sumTemperature += temperature
			    measurement.countTemperature++
			}
			
			data[location] = measurement
		}

		locations := make([]string, 0, len(data))
		for name := range data {
			locations = append(locations, name)
		}
		
		sort.Strings(locations)

		fmt.Printf("{")
		for _, name := range locations {
			  measurement := data[name]
			  fmt.Printf(
				    "%s=%.1f/%.1f/%.1f, ",
					  name,
					  measurement.minTemperature,
					  measurement.sumTemperature/float64(measurement.countTemperature),
					  measurement.maxTemperature,
			  )
		}
		fmt.Printf("}\n")
		fmt.Println(time.Since(start))
}