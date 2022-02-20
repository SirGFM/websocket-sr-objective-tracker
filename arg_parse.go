package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Args struct {
	// IP on which the server will accept connections. Defaults to 127.0.0.1.
	IP string
	// Port on which the server will accept connections. Defaults to 8000.
	Port int
	// Resource directory where files may be read from. Defaults to 'res/'.
	ResDir string
}

// parseAndMergeFile reads the given file and merge it with the previously
// parsed command-line arguments.
func parseAndMergeFile(args *Args, def Args, f *os.File) {
	var jsonArgs Args

	dec := json.NewDecoder(f)
	err := dec.Decode(&jsonArgs)
	if err != nil {
		log.Fatalf("Couldn't decode the configuration file '%+v': %+v", f.Name(), err)
	}

	// NOTE: There's no reasonable way to differentiate between a value
	// that wasn't set and the empty/zero value (no, I'm not using pointers
	// for this). So, if a value isn't set, revert it to the default value.
	if len(jsonArgs.IP) == 0 {
		jsonArgs.IP = def.IP
	}
	if jsonArgs.Port == 0 {
		jsonArgs.Port = def.Port
	}
	if len(jsonArgs.ResDir) == 0 {
		jsonArgs.ResDir = def.ResDir
	}

	// Walk over every set argument to override the JSON file
	flag.Visit(func (f *flag.Flag) {
		if f.Name == "confFile" {
			// Skip the file itself
			return
		}

		var tmp interface{}
		tmp = f.Value
		get, ok := tmp.(flag.Getter)
		if !ok {
			log.Fatalf("'%s' doesn't have an associated flag.Getter", f.Name)
		}

		switch f.Name {
		case "IP":
			val, _ := get.Get().(string)
			log.Printf("Overriding JSON's IP (%+v) with CLI's value (%+v)", jsonArgs.IP, val)
			jsonArgs.IP = val
		case "Port":
			val, _ := get.Get().(int)
			log.Printf("Overriding JSON's Port (%+v) with CLI's value (%+v)", jsonArgs.Port, val)
			jsonArgs.Port = val
		case "ResDir":
			val, _ := get.Get().(string)
			log.Printf("Overriding JSON's ResDir (%+v) with CLI's value (%+v)", jsonArgs.ResDir, val)
			jsonArgs.ResDir = val
		}
	})

	*args = jsonArgs
}

// parseArgs either from the command line or from the supplied JSON file.
//
// If a JSON file is supplied, it's used as the default parameters, which may be overriden by CLI-supplied arguments.
func parseArgs() Args {
	var args Args
	var confFile string
	defArgs := Args {
		IP: "127.0.0.1",
		Port: 8000,
		ResDir: "res/",
	}
	const defaultConfFile = "config.json"

	flag.StringVar(&args.IP, "IP", defArgs.IP, "IP on which the server will accept connections")
	flag.IntVar(&args.Port, "Port", defArgs.Port, "Port on which the server will accept connections")
	flag.StringVar(&args.ResDir, "ResDir", defArgs.ResDir, "Resource directory where files may be read from")
	flag.StringVar(&confFile, "confFile", defaultConfFile, "JSON file with the configuration options. May be overriden by other CLI arguments")
	flag.Parse()

	if len(confFile) != 0 {
		f, err := os.Open(confFile)
		if err != nil && confFile == defaultConfFile {
			// Ignore errors if trying to read the default file.
		} else if err != nil {
			log.Fatalf("Couldn't open the configuration file '%+v': %+v", confFile, err)
		} else {
			defer f.Close()
			parseAndMergeFile(&args, defArgs, f)
		}
	}

	log.Printf("Starting server with options:")
	log.Printf("  - IP: %+v", args.IP)
	log.Printf("  - Port: %+v", args.Port)
	log.Printf("  - ResDir: %+v", args.ResDir)

	return args
}
