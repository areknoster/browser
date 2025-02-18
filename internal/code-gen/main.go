// This is part of an internal code generation tool of Gost-DOM. It's in an
// internal package and not used by production code.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	htmlelements "github.com/gost-dom/code-gen/html-elements"
	wrappers "github.com/gost-dom/code-gen/script-wrappers"
)

func getWriter(output string) io.Writer {
	if output == "stdout" {
		return os.Stdout
	}
	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating output file")
		os.Exit(1)
	}
	return file
}

var generators = map[string]func(io.Writer) error{
	"html-elements": generateHtmlElements,
}

func main() {
	debug := flag.Bool("d", false, "Debug")
	outputFile := flag.String("o", "", "Output file to write")
	generatorType := flag.String("g", "", "Generator type")
	packageName := flag.String("p", "", "Package to generate")
	flag.Parse()
	switch *generatorType {
	case "goja":
		gen := wrappers.NewGojaWrapperModuleGenerator()
		exitOnError(gen.GenerateScriptWrappers())
		os.Exit(0)
		return
	case "scripting":
		gen := wrappers.NewScriptWrapperModulesGenerator()
		exitOnError(gen.GenerateScriptWrappers())
		os.Exit(0)
		return
	case "htmlelements":
		exitOnError(htmlelements.CreateImplementationPackage("html"))
		os.Exit(0)
		return
	case "dom":
		exitOnError(htmlelements.CreateImplementationPackage("dom"))
		os.Exit(0)
		return
	case "gotypes":
		if packageName == nil {
			panic("Missing package spec")
		}
		exitOnError(htmlelements.CreateImplementationPackage(*packageName))
		os.Exit(0)
	}

	if *outputFile == "" || *generatorType == "" {
		fmt.Println("Internal code generator from IDL definitions")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *debug {
		fmt.Println(strings.Join(os.Args, " "))
		fmt.Println("--------")
	}

	file := getWriter(*outputFile)

	generator, ok := generators[*generatorType]
	if !ok {
		os.Exit(1)
	}
	err := generator(file)
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("Error running generator")
		fmt.Println(err)
		os.Exit(1)
	}
}
