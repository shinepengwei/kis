package main

import(
	"fmt"
	"kis/apiserver"
	"os"
)
func main(){
	command:=apiserver.NewAPIServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}