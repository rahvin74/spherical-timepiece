package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rahvin74/spherical-timepiece/internal/timepiece"
	"github.com/rahvin74/spherical-timepiece/internal/timesphere"
)

func main() {
	fmt.Println("Welcome To")
	bannerText := banner()
	fmt.Println(bannerText)
	fmt.Println()

	fmt.Println("This program mimics a Ball Clock to calculate how many days")
	fmt.Println("it would take for balls running through the clock to return")
	fmt.Println("to their original starting positions.")
	fmt.Println()
	fmt.Println("Enjoy!")

	for {
		var totalBalls string

		fmt.Println()
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Println("Enter the desired number of balls the clock will run.")
		fmt.Printf("Please choose a number between 27 - 127 (q to quit): ")

		fmt.Scanln(&totalBalls)

		if totalBalls == "q" || totalBalls == "Q" {
			fmt.Println("Thank you for using Spherical Timepiece! May all your times be wonderful!")
			fmt.Println()
			break
		}

		totalBallsInt, err := strconv.Atoi(totalBalls)
		if err != nil {
			fmt.Println("You must enter a number.")
			continue
		}

		spheres := make([]timesphere.MinuteBall, totalBallsInt, totalBallsInt)

		for i := range spheres {
			spheres[i].SetOriginalPosition(i)
		}

		clock := timepiece.NewMechanism()

		t := time.Now()
		runTime := clock.Run(spheres)
		totalTime := time.Since(t)

		fmt.Println()
		fmt.Printf("%d balls will take %d days. Processing time: %d Milliseconds\n", len(spheres), runTime, totalTime.Milliseconds())
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Println()
	}
}

func banner() string {
	bannerText := `
   _____       __              _            __   _______                      _              
  / ___/____  / /_  ___  _____(_)________ _/ /  /_  __(_)___ ___  ___  ____  (_)__  ________ 
  \__ \/ __ \/ __ \/ _ \/ ___/ / ___/ __ '/ /    / / / / __ '__ \/ _ \/ __ \/ / _ \/ ___/ _ \
 ___/ / /_/ / / / /  __/ /  / / /__/ /_/ / /    / / / / / / / / /  __/ /_/ / /  __/ /__/  __/
/____/ .___/_/ /_/\___/_/  /_/\___/\__,_/_/    /_/ /_/_/ /_/ /_/\___/ .___/_/\___/\___/\___/ 
    /_/                                                            /_/                       
	`
	return bannerText
}
