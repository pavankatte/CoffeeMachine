package main

import (
    "fmt"
    "os"
    "io/ioutil"
)


func main() {
    //this main function will take input JSON file as command line argument and invokes ExecuteCoffeeMachine
    if len(os.Args) < 2 {
        fmt.Println("Missing parameter, provide file name!")
        return
    }
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println("Can't read file:", os.Args[1])
        panic(err)
    }

    ExecuteCoffeeMachine(data)
}