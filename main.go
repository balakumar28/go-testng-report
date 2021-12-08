package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	cmd := flag.NewFlagSet("testng-report", flag.ExitOnError)
	json := cmd.String("json-report", "report.json", "Golang json test report")
	groups := cmd.String("testng-groups", "groups.properties", "TestNG group mapping")
	report := cmd.String("out", "testng-results.xml", "TestNG Report XML")
	_ = cmd.Parse(os.Args[1:])

	err := GenerateReport(*json, *groups, *report)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Usage:\tgo-testng-report [OPTIONS]\n\n")
}