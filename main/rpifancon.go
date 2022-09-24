/*
A very simple fan controller for a raspberry pi.
It reads temperatures over a configurable range, gets the mean
and compares that to a pre-set threshhold. If its above that it turns
a fan on via a GPIO pin for a few seconds and re-tests.
*/

package main

import (
	"flag"
	"fmt"
	rpio "github.com/stianeikeland/go-rpio/v4"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func currtemp() int {
	out, err := exec.Command("vcgencmd", "measure_temp").Output()

	if err != nil {
		log.Fatal(err)
	}

	currentTemp := string(out)

	// what? working outwards: replace 'C with nothing, then replace temp= with nothing, grab the first element which should be
	// the only one, convert this string to an int. This should give us the temp in C rounded down to the nearest int
	temp, err := strconv.Atoi(strings.Split(strings.Replace(strings.Replace(currentTemp, "'C", "", -1), "temp=", "", -1), ".")[0])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	return temp
}

func getmean(tempValues []int) int {

	// get length of array
	arrSize := len(tempValues)

	sum := 0

	for i := 0; i < arrSize; i++ {
		sum += (tempValues[i])
	}

	// avg to find the average
	avg := int((float64(sum)) / (float64(arrSize)))

	return avg
}

func main() {

	// get our cli options
	iopin := flag.Int("iopin", 17, " The GPIO pin used to control the fan.")
	threshold := flag.Int("threshold", 65, " The temperature in celius above which to enable the fan.")
	debug := flag.Bool("debug", false, " Debug info.")
	wait := flag.Int("wait", 5, " The amount of time to wait between polling temperature. Multiply this by checks to get time between pin state changes.")
	checks := flag.Int("checks", 3, " The number of checks of temperatures to check before state changes.")
	help := flag.Bool("help", false, " This info.")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// set pin
	pin := rpio.Pin(*iopin)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	// build our array of recent temperatures and populate it initially
	var recenttemps []int

	for x := 0; x < *checks; x++ {

		if *debug {
			log.Printf("Doing Initial Temp Reads...")
			log.Printf("Inital temp read %i:  %s", x, currtemp())
		}

		recenttemps = append(recenttemps, currtemp())

		// pin.Toggle()
		time.Sleep(time.Second * time.Duration(*wait))
	}

	if *debug {
		log.Printf("Inital temp reads:  %a", recenttemps)
	}

	// keep watching the temperature, and if it goes above our threshold turn the fan on
	for {

		recenttemps = append(recenttemps, currtemp())

		_, recenttemps = recenttemps[0], recenttemps[1:]

		if *debug {
			log.Printf("Temp reads:  %a", recenttemps)
			log.Printf("Mean temp:  %i", getmean(recenttemps))
		}

		time.Sleep(time.Second * time.Duration(*wait))

		avg := getmean(recenttemps)

		if *debug {
			if rpio.High == pin.Read() {
				log.Printf("Pin High")
			} else if rpio.Low == pin.Read() {
				log.Printf("Pin Low")
			}
		}

		if (avg >= *threshold) && (rpio.Low == pin.Read()) {
			// turn fan on
			log.Printf("Turning fan on.")
			pin.High()
			// sleep for a bit before re-entering the loop to delay any possible state changes
			time.Sleep((time.Second * time.Duration(*wait)) * time.Duration(*checks))
		} else if (avg < *threshold) && rpio.High == pin.Read() {
			// turn fan off
			log.Printf("Turning fan off.")
			pin.Low()
			// sleep for a bit before re-entering the loop to delay any possible state changes
			time.Sleep((time.Second * time.Duration(*wait)) * time.Duration(*checks))
		} else {
			// No changes to make
			if *debug {
				log.Printf("Pin in correct state...")
			}
		}

	}

}
