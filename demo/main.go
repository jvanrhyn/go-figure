package main

import "github.com/jvanrhyn/go-figure"

func main() {
	myFigure, err := figure.NewFigure("Hello World", "", true)
	if err != nil {
		panic(err)
	}
	myFigure.Print()
}
