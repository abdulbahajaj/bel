package main

import (
  "io/ioutil"
  "log"
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

  module, err := reader.Read(program)
  if err != nil{
    if err.Error() != "EmptyTokenList"{
      log.Fatal(err)
      return
    }
  }

  env := primitives.GetPrimitiveEnv()
  eval.Eval(module, env)
}
