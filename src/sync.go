package main

import "fmt"

func runSync(steps tSyncSteps) {
	for _, step := range steps {
		fmt.Printf("%+v\n", step)
	}
}
