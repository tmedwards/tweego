/*
	option (a simple command-line option parser for Go)

	Copyright © 2014–2018 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

// Package option implements simple command-line option parsing.
package option

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// OptionTerminator is the string, when seen on the command line, which terminates further option processing.
const OptionTerminator = "--"

// OptionTypeMap is the map of recognized type abbreviations to types.
var OptionTypeMap = map[string]string{"s": "string", "i": "int", "u": "uint", "f": "float", "b": "bool"}

// Config is {TODO}.
type Config struct {
	Name       string
	Definition string
	Flags      int
	//Default interface{}
}

// Options is {TODO}.
type Options struct {
	Definitions []Config
}

/*
// NewOptions is {TODO}.
func NewOptions(options ...Config) Options {
	return Options{options}
}
*/

// NewParser returns {TODO}.
func NewParser() Options {
	return Options{}
}

// Add adds a new option definition.
func (optDef *Options) Add(name, def string /*, flags int*/) {
	optDef.Definitions = append(optDef.Definitions, Config{name, def, 0 /*flags*/})
}

type optionDefinition struct {
	name       string
	wantsValue bool
	valueType  string
	repeatable bool
	flags      int
}
type optionMap map[string]optionDefinition

func (optDef Options) buildOptionMap() optionMap {
	optMap := make(optionMap)
	for _, def := range optDef.Definitions {
		if def.Definition != "" {
			names, opts := parseDefinition(def.Definition)
			for i := range names {
				opts[i].name = def.Name
				opts[i].flags = def.Flags
				optMap[names[i]] = opts[i]
			}
		}
	}
	return optMap
}

func parseDefinition(optSpec string) ([]string, []optionDefinition) {
	var (
		names []string
		defs  []optionDefinition
	)
	for _, def := range strings.Split(optSpec, "|") {
		if i := strings.LastIndex(def, "="); i != -1 {
			// value receiving option
			names = append(names, def[:i])
			optDef := optionDefinition{wantsValue: true}
			valueType := def[i+1:]
			if valueType == "s+" || valueType == "i+" || valueType == "u+" || valueType == "f+" {
				// special case: value receiving + repeatable
				optDef.repeatable = true
				optDef.valueType = OptionTypeMap[valueType[:1]]
			} else if _, ok := OptionTypeMap[valueType]; ok {
				// normal cases
				optDef.valueType = OptionTypeMap[valueType]
			} else {
				// what type now?
				panic(fmt.Errorf("Cannot parse value type %q in option specification %q.", valueType, optSpec))
			}
			defs = append(defs, optDef)
		} else if i := strings.LastIndex(def, "+"); i != -1 {
			// repeatable unsigned integer option
			names = append(names, def[:i])
			defs = append(defs, optionDefinition{
				repeatable: true,
				valueType:  OptionTypeMap["u"],
			})
		} else {
			// void/empty option
			names = append(names, def)
			defs = append(defs, optionDefinition{})
		}
	}
	return names, defs
}

// ParsedOptionsMap is {TODO}.
type ParsedOptionsMap map[string]interface{}

// ParseCommandLine returns {TODO}.
func (optDef Options) ParseCommandLine() (ParsedOptionsMap, []string, error) {
	return optDef.Parse(os.Args[1:])
}

// Parse returns {TODO}.
func (optDef Options) Parse(args []string) (ParsedOptionsMap, []string, error) {
	var (
		passThrough []string
		err         error
	)
	options := make(ParsedOptionsMap)

	optMap := optDef.buildOptionMap()

	for i, argc := 0, len(args); i < argc; i++ {
		var (
			name string
		)
		sz := len(args[i])
		if sz > 1 && args[i][0] == '-' {
			// could be an option, try to parse it
			if eqPos := strings.Index(args[i], "="); eqPos != -1 {
				// with bundled value
				name = args[i][:eqPos]
				if opt, ok := optMap[name]; ok {
					if opt.wantsValue {
						if value, err := convertType(args[i][eqPos+1:], opt.valueType); err == nil {
							if opt.repeatable {
								if _, ok := options[opt.name]; !ok {
									switch opt.valueType {
									case "string":
										options[opt.name] = make([]string, 0, 4)
									case "int":
										options[opt.name] = make([]int, 0, 4)
									case "uint":
										options[opt.name] = make([]uint, 0, 4)
									case "float":
										options[opt.name] = make([]float64, 0, 4)
									}
								}
								switch opt.valueType {
								case "string":
									options[opt.name] = append(options[opt.name].([]string), value.(string))
								case "int":
									options[opt.name] = append(options[opt.name].([]int), value.(int))
								case "uint":
									options[opt.name] = append(options[opt.name].([]uint), value.(uint))
								case "float":
									options[opt.name] = append(options[opt.name].([]float64), value.(float64))
								}
							} else {
								options[opt.name] = value
							}
						} else {
							err = fmt.Errorf("Option %q %s.", name, err.Error())
							break
						}
					} else {
						err = fmt.Errorf("Option %q does not take a value.", name)
						break
					}
				} else {
					err = fmt.Errorf("Unknown option %q.", name)
					break
				}
			} else {
				// without bundled value
				name = args[i]
				if name == OptionTerminator {
					// processing terminated, pass any remaining arguments on through
					passThrough = append(passThrough, args[i+1:]...)
					break
				}
				if opt, ok := optMap[name]; ok {
					if opt.wantsValue {
						i++
						if i < argc {
							if value, err := convertType(args[i], opt.valueType); err == nil {
								if opt.repeatable {
									if _, ok := options[opt.name]; !ok {
										switch opt.valueType {
										case "string":
											options[opt.name] = make([]string, 0, 4)
										case "int":
											options[opt.name] = make([]int, 0, 4)
										case "uint":
											options[opt.name] = make([]uint, 0, 4)
										case "float":
											options[opt.name] = make([]float64, 0, 4)
										}
									}
									switch opt.valueType {
									case "string":
										options[opt.name] = append(options[opt.name].([]string), value.(string))
									case "int":
										options[opt.name] = append(options[opt.name].([]int), value.(int))
									case "uint":
										options[opt.name] = append(options[opt.name].([]uint), value.(uint))
									case "float":
										options[opt.name] = append(options[opt.name].([]float64), value.(float64))
									}
								} else {
									options[opt.name] = value
								}
							} else {
								err = fmt.Errorf("Option %q %s.", name, err.Error())
								break
							}
						} else {
							err = fmt.Errorf("Option %q requires a value.", name)
							break
						}
					} else if opt.repeatable {
						if _, ok := options[opt.name]; ok {
							options[opt.name] = options[opt.name].(uint) + 1
						} else {
							options[opt.name] = 1
						}
					} else {
						options[opt.name] = true
					}
				} else {
					err = fmt.Errorf("Unknown option %q.", name)
					break
				}
			}
		} else {
			// not an option, pass it through
			passThrough = append(passThrough, args[i])
		}
	}
	return options, passThrough, err
}

func convertType(original, targetType string) (interface{}, error) {
	var (
		value interface{}
		err   error
	)
	switch targetType {

	case "string":
		value = original

	case "int":
		var tmp int64
		if tmp, err = strconv.ParseInt(original, 10, 0); err != nil {
			err = fmt.Errorf("Cannot interpret value %q as an integer: %s.", original, err.Error())
			break
		}
		value = int(tmp)

	case "uint":
		var tmp uint64
		if tmp, err = strconv.ParseUint(original, 10, 0); err != nil {
			err = fmt.Errorf("Cannot interpret value %q as an unsigned integer: %s.", original, err.Error())
			break
		}
		value = uint(tmp)

	case "float":
		var tmp float64
		if tmp, err = strconv.ParseFloat(original, 64); err != nil {
			err = fmt.Errorf("Cannot interpret value %q as a floating-point number: %s.", original, err.Error())
			break
		}
		value = tmp

	case "bool":
		var tmp bool
		if tmp, err = strconv.ParseBool(original); err != nil {
			err = fmt.Errorf("Cannot interpret value %q as a boolean: %s.", original, err.Error())
			break
		}
		value = bool(tmp)

	}
	return value, err
}
