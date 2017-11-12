package main

import (
	"flag"
	"log"

	"ciel/internal/cieldir.1"
)

func addInst() {
	basePath := flagCielDir()
	instName := flagInstance()
	parse()
	if *instName == "" {
		*instName = flag.Arg(0)
	}

	i := &cieldir.CielDir{BasePath: *basePath}
	i.Check()

	if i.InstExists(*instName) {
		log.Fatalln("already has " + *instName)
	}
	i.InstAdd(*instName)
	i.InstMount(*instName)
}

func delInst() {
	basePath := flagCielDir()
	instName := flagInstance()
	parse()
	if *instName == "" {
		*instName = flag.Arg(0)
	}

	i := &cieldir.CielDir{BasePath: *basePath}
	i.Check()
	i.CheckInst(*instName)

	i.InstUnmount(*instName)
	i.InstDel(*instName)
}