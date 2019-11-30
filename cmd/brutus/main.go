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

  fmt.Println("Program:")
  fmt.Println(program)
  module, err := reader.Read(program)
  if err != nil{
    fmt.Println(err)
  }

  // fmt.Println(len(module.Expressions))
  // fmt.Println(module.Expressions)

  fmt.Println("Parsed:")
  reader.PrintModule(module)


  fmt.Println("Result: ")
  env := primitives.GetPrimitiveEnv()
  return_stack, _ := eval.Eval(module, env)
  fmt.Println(return_stack)
}
