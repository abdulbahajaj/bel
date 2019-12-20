package main

import (
  "io/ioutil"
  "fmt"
  "os"

  "github.com/abdulbahajaj/brutus/pkg/reader"
  "github.com/abdulbahajaj/brutus/pkg/eval"
  "github.com/abdulbahajaj/brutus/pkg/primitives"
)

func main() {

  args := os.Args[1:]
  fileName := args[0]

  dat, err := ioutil.ReadFile(fileName)
  if err != nil {
    panic(err)
  }
  program := string(dat)
  // fmt.Println(string(program[3]) == "\n")
  // fmt.Println("Program:")
  // fmt.Println(program)

  module, err := reader.Read(program)
  if err != nil{
    fmt.Println(err)
  }

  // fmt.Println(">Parsed:")
  // reader.PrintModule(module)


  // fmt.Println(">Started eval: ")
  env := primitives.GetPrimitiveEnv()
  eval.Eval(module, env)
  // fmt.Println(">Result: ")
  // fmt.Println(return_stack)
}
