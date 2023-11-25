package main // "github.com/privacy-pal/privacy-pal/go/cmd/genpal"

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"log"
	"os"
	"strings"

	genpal "github.com/privacy-pal/privacy-pal/go/internal/genpal"
	"golang.org/x/tools/go/packages"
	yaml "gopkg.in/yaml.v2"
)

var (
	mode   = flag.String("mode", "", "typenames or yamlspec")
	input  = flag.String("input", "", "if mode is typenames, comma-separated list of type names; if mode is yamlspec, path to yaml file")
	output = flag.String("output", "", "output file name; default srcdir/<type>_privacy.go")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of genpal:\n")
	fmt.Fprintf(os.Stderr, "\tgenpal [flags]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func validateArgs() (genpal.Mode, error) {
	// Validate mode
	switch *mode {
	case string(genpal.ModeTypenames):
		types := strings.Split(*input, ",")
		if len(types) == 0 {
			return "", fmt.Errorf("no typenames provided")
		}
	case string(genpal.ModeYamlspec):
		if *input == "" {
			return "", fmt.Errorf("no yamlspec provided")
		}
		// check if file exists
		_, err := os.Stat(*input)
		if err != nil {
			return "", fmt.Errorf("yamlspec file does not exist")
		}
	case "":
		return "", fmt.Errorf("no mode provided")
	default:
		return "", fmt.Errorf("invalid mode")
	}
	return genpal.Mode(*mode), nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("genpal: ")
	flag.Usage = Usage
	flag.Parse()

	genpalMode, err := validateArgs()
	if err != nil {
		log.Printf("%s. See 'genpal -help'.", err)
		os.Exit(2)
	}

	outputName := *output
	var args []string

	if outputName == "" {
		// Write to a file under current working directory if no output specified
		args = []string{"."}
		outputName = "./privacy_genpal.go"
	} else {
		// set args to be the directory containing output file
		if !(strings.HasPrefix(outputName, "/") || strings.HasPrefix(outputName, "./") || strings.HasPrefix(outputName, "../")) {
			outputName = "./" + outputName
		}
		args = []string{outputName[:strings.LastIndex(outputName, "/")]}
	}

	// Parse the package once.
	g := Generator{}
	g.parsePackage(args)

	// Print the header and package clause.
	g.Printf("// Code generated by \"genpal %s\"\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")
	g.Printf("import (\npal \"github.com/privacy-pal/privacy-pal/go/pkg\"\n)\n\n")

	if genpalMode == genpal.ModeTypenames {
		types := strings.Split(*input, ",")
		g.Printf(genpal.GenerateWithTypenameMode(types))

	} else if genpalMode == genpal.ModeYamlspec {
		data, err := os.ReadFile(*input)
		if err != nil {
			log.Printf("error reading yamlspec file: %s\n", err)
			os.Exit(2)
		}
		var dataNodes map[string]genpal.DataNodeProperty
		err = yaml.Unmarshal(data, &dataNodes)
		if err != nil {
			log.Printf("error unmarshalling yamlspec file: %s\n", err)
			os.Exit(2)
		}

		var mapSlice yaml.MapSlice
		err = yaml.Unmarshal(data, &mapSlice)
		if err != nil {
			log.Printf("error unmarshalling yamlspec file: %s\n", err)
			os.Exit(2)
		}
		typenames := make([]string, 0)
		for _, item := range mapSlice {
			typenames = append(typenames, item.Key.(string))
		}
		g.Printf(genpal.GenerateWithYamlspecMode(typenames, dataNodes))
	}

	// Format the output.
	src := g.format()

	// Write to file.
	err = os.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.

	logf func(format string, args ...interface{}) // test logging hook; nil when not testing
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
}

type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests: false,
		Logf:  g.logf,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages matching %v", len(pkgs), strings.Join(patterns, " "))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}