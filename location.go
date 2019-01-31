package main

import (
	"fmt"

	"./generators"
)

func main() {
	adjective, err := generators.RandomAdjective()
	errorHandler(err)
	fmt.Println(adjective)
	things := generators.AssembleBuildings()
	tmp := things[generators.Hospitality][0]
	fmt.Println(tmp)
	fmt.Println(tmp.Weight)

}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
