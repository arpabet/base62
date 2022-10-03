/**
  Copyright (c) 2022 Zander Schwid & Co. LLC. All rights reserved.
*/

package app

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/schwid/base62"
	"io"
	"os"
	"runtime"
	"unicode"

	"github.com/jessevdk/go-flags"
)

type app struct {

	name string
	version string
	build string

	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

type flagopts struct {
	Decode   bool             `short:"D" long:"decode" description:"decodes input"`
	Input    []string         `short:"i" long:"input" default:"-" description:"input file"`
	Output   string           `short:"o" long:"output" default:"-" description:"output file"`
	Version  bool             `short:"v" long:"version" description:"print version"`
}

func Run(name, version, build  string) error {
	return (&app{
		name: name,
		version: version,
		build: build,
		inStream:  os.Stdin,
		outStream: os.Stdout,
		errStream: os.Stderr,
	}).run(os.Args[1:])
}

func (cli *app) run(args []string) error {
	var opts flagopts
	args, err := flags.NewParser(
		&opts, flags.HelpFlag|flags.PassDoubleDash,
	).ParseArgs(args)
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			fmt.Fprintln(cli.outStream, err.Error())
			return nil
		}
		return err
	}
	if opts.Version {
		fmt.Fprintf(cli.outStream, "%s %s (build: %s/%s)\n", cli.name, cli.version, cli.build, runtime.Version())
		return nil
	}
	var inputFiles []string
	for _, name := range append(opts.Input, args...) {
		if name != "" && name != "-" {
			inputFiles = append(inputFiles, name)
		}
	}
	if opts.Output != "-" {
		file, err := os.Create(opts.Output)
		if err != nil {
			return err
		}
		defer file.Close()
		cli.outStream = file
	}
	var result error
	if len(inputFiles) == 0 {
		if err := cli.runInternal(opts.Decode, cli.inStream); err != nil {
			result = err
		}
	}
	for _, name := range inputFiles {
		if err := cli.runFile(opts.Decode, name); err != nil {
			result = err
		}
	}
	return result
}

func (cli *app) runFile(decode bool, name string) error {
	file, err := os.Open(name)
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		return err
	}
	defer file.Close()
	return cli.runInternal(decode, file)
}

func (cli *app) runInternal(decode bool, in io.Reader) error {
	scanner := bufio.NewScanner(in)
	var status error
	var result []byte
	var err error
	for scanner.Scan() {
		src := scanner.Bytes()
		if decode {
			result, err = processLine(src, func(in []byte) ([]byte, error) {
				return base62.StdEncoding.DecodeString(string(in))
			})
		} else {
			result, err = processLine(src, func(in []byte) ([]byte, error) {
				return []byte(base62.StdEncoding.EncodeToString(in)), nil
			})
		}
		if err != nil {
			fmt.Fprintln(cli.errStream, err.Error()) // should print error each line
			status = err
			continue
		}
		cli.outStream.Write(result)
		cli.outStream.Write([]byte{0x0a})
	}
	return status
}

func processLine(src []byte, f func([]byte) ([]byte, error)) ([]byte, error) {
	var i, j int
	var res []byte
	for j < len(src) {
		j = bytes.IndexFunc(src[i:], unicode.IsSpace)
		if j >= 0 {
			j += i
		} else {
			j = len(src)
		}
		got, err := f(src[i:j])
		if err != nil {
			return nil, err
		}
		res = append(res, got...)
		if j == len(src) {
			break
		}
		i = bytes.IndexFunc(src[j:], func(r rune) bool { return !unicode.IsSpace(r) })
		if i >= 0 {
			i += j
		} else {
			i = len(src)
		}
		res = append(res, src[j:i]...)
		if i == len(src) {
			break
		}
	}
	return res, nil
}


