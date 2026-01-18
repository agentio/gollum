package lexica

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/iancoleman/strcase"
)

func (lexicon *Lexicon) generateLeafCommands(root string) {
	allow := []string{
		"com.atproto.admin.getInviteCodes",
		"com.atproto.admin.updateAccountPassword",
		"com.atproto.server.createInviteCode",
		"com.atproto.server.getAccountInviteCodes",
		"com.atproto.sync.listRepos",
	}
	if !slices.Contains(allow, lexicon.Id) {
		//return
	}
	for defname, def := range lexicon.Defs {
		switch def.Type {
		case "query":
			lexicon.generateLeafCommandForDef(root, defname, def)
		case "procedure":
			lexicon.generateLeafCommandForDef(root, defname, def)
		}
	}
}

func (lexicon *Lexicon) generateLeafCommandForDef(root, defname string, def *Def) {
	if defname != "main" {
		log.Errorf("Can't generate leaf command for %s %s", lexicon.Id, defname)
	}
	id := strings.Replace(lexicon.Id, ".", "-", 2) // merge initial segments of the lexicon id
	dirname := strings.ToLower(root + "/" + strings.ReplaceAll(id, ".", "/"))
	os.MkdirAll(dirname, 0755)
	filename := dirname + "/cmd.go"

	parts := strings.Split(id, ".")
	lastpart := parts[len(parts)-1]
	packagename := strings.ToLower(lastpart)
	commandname := strcase.ToKebab(lastpart)
	handlerName := idPrefix(lexicon.Id)

	s := &strings.Builder{}
	fmt.Fprintf(s, "package %s // %s\n\n", packagename, lexicon.Id)
	fmt.Fprintf(s, "import \"github.com/spf13/cobra\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/gen/xrpc\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/common\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/client\"\n")
	fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
	fmt.Fprintf(s, "var _output string\n")
	if def.Type == "query" && def.Parameters != nil {
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyValue := def.Parameters.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "var %s string\n", propertyName)
			case "integer":
				fmt.Fprintf(s, "var %s int64\n", propertyName)
			case "boolean":
				fmt.Fprintf(s, "var %s bool\n", propertyName)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "var %s []string\n", propertyName)
				} else {
					fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyName, propertyValue)
				}
			default:
				fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyName, propertyValue)
			}
		}
	} else if def.Type == "procedure" && def.Input != nil {
		for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
			propertyValue := def.Input.Schema.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "var %s string\n", propertyName)
			case "integer":
				fmt.Fprintf(s, "var %s int64\n", propertyName)
			case "boolean":
				fmt.Fprintf(s, "var %s bool\n", propertyName)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "var %s []string\n", propertyName)
				} else {
					fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyName, propertyValue)
				}
			case "unknown":
				fmt.Fprintf(s, "var %s string\n", propertyName)
			default:
				fmt.Fprintf(s, "// FIXME var %s %+v\n", propertyName, propertyValue)
			}
		}
	}
	fmt.Fprintf(s, "cmd := &cobra.Command{\n")
	fmt.Fprintf(s, "Use: \"%s\",\n", commandname)
	fmt.Fprintf(s, "Short: common.Truncate(xrpc.%s_Description),\n", handlerName)
	fmt.Fprintf(s, "Long: xrpc.%s_Description,\n", handlerName)
	fmt.Fprintf(s, "Args: cobra.NoArgs,\n")
	fmt.Fprintf(s, "RunE: func(cmd *cobra.Command, args []string) error {\n")
	if def.Type == "query" {
		fmt.Fprintf(s, "client := client.NewClient()\n")
		fmt.Fprintf(s, "response, err := xrpc.%s(\n", handlerName)
		fmt.Fprintf(s, "cmd.Context(),\n")
		fmt.Fprintf(s, "client,\n")
		if def.Parameters != nil {
			for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
				propertyValue := def.Parameters.Properties[propertyName]
				switch propertyValue.Type {
				case "string":
					fmt.Fprintf(s, "%s,\n", propertyName)
				case "integer":
					fmt.Fprintf(s, "%s,\n", propertyName)
				case "boolean":
					fmt.Fprintf(s, "%s,\n", propertyName)
				case "array":
					if propertyValue.Items.Type == "string" {
						fmt.Fprintf(s, "%s,\n", propertyName)
					}
				default:
				}
			}
		}
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "if err != nil {return err}\n")
		fmt.Fprintf(s, "return common.Write(cmd.OutOrStdout(), _output, response)\n")
	} else if def.Type == "procedure" && (def.Input == nil || def.Input.Encoding == "application/json") {
		fmt.Fprintf(s, "client := client.NewClient()\n")
		resultIfNeeded := ""
		if def.Output != nil {
			resultIfNeeded = "response, "
		}
		fmt.Fprintf(s, "%serr := xrpc.%s(\n", resultIfNeeded, handlerName)
		fmt.Fprintf(s, "cmd.Context(),\n")
		fmt.Fprintf(s, "client,\n")
		if def.Input != nil {
			fmt.Fprintf(s, "&xrpc.%s_Input{\n", handlerName)
			for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
				propertyValue := def.Input.Schema.Properties[propertyName]
				switch propertyValue.Type {
				case "string":
					if !slices.Contains(def.Input.Schema.Required, propertyName) {
						fmt.Fprintf(s, "%s: common.StringPointerOrNil(%s),\n", capitalize(propertyName), propertyName)
					} else {
						fmt.Fprintf(s, "%s: %s,\n", capitalize(propertyName), propertyName)
					}
				case "integer":
					if !slices.Contains(def.Input.Schema.Required, propertyName) {
						fmt.Fprintf(s, "%s: common.Int64PointerOrNil(%s),\n", capitalize(propertyName), propertyName)
					} else {
						fmt.Fprintf(s, "%s: %s,\n", capitalize(propertyName), propertyName)
					}
				case "boolean":
					//fmt.Fprintf(s, "%s,\n", propertyName)
				case "array":
					if propertyValue.Items.Type == "string" {
						//	fmt.Fprintf(s, "%s,\n", propertyName)
					}
				default:
				}
			}
			fmt.Fprintf(s, "},\n")
		}
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "if err != nil {return err}\n")
		if def.Output == nil {
			fmt.Fprintf(s, "return nil\n")
		} else {
			fmt.Fprintf(s, "return common.Write(cmd.OutOrStdout(), _output, response)\n")
		}
	} else {
		fmt.Fprintf(s, "return errors.New(\"unimplemented\")")
	}
	fmt.Fprintf(s, "},\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "cmd.Flags().StringVarP(&_output, \"output\", \"o\", \"\", \"output destination\")\n")
	if def.Type == "query" && def.Parameters != nil {
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyValue := def.Parameters.Properties[propertyName]
			flagName := strcase.ToKebab(propertyName)
			description := propertyName
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s\")\n", propertyName, flagName, description)
			case "integer":
				fmt.Fprintf(s, "cmd.Flags().Int64Var(&%s, \"%s\", %d, \"%s\")\n", propertyName, flagName, int64Value(propertyValue.Default), description)
			case "boolean":
				fmt.Fprintf(s, "cmd.Flags().BoolVar(&%s, \"%s\", %t, \"%s\")\n", propertyName, flagName, boolValue(propertyValue.Default), description)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "cmd.Flags().StringArrayVar(&%s, \"%s\", nil, \"%s\")\n", propertyName, flagName, description)
				} else {
					fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyName, propertyValue)
				}
			default:
				fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyName, propertyValue)
			}
		}
	} else if def.Type == "procedure" && def.Input != nil {
		for _, propertyName := range sortedPropertyNames(def.Input.Schema.Properties) {
			propertyValue := def.Input.Schema.Properties[propertyName]
			flagName := strcase.ToKebab(propertyName)
			description := propertyName
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s\")\n", propertyName, flagName, description)
			case "integer":
				fmt.Fprintf(s, "cmd.Flags().Int64Var(&%s, \"%s\", %d, \"%s\")\n", propertyName, flagName, int64Value(propertyValue.Default), description)
			case "boolean":
				fmt.Fprintf(s, "cmd.Flags().BoolVar(&%s, \"%s\", %t, \"%s\")\n", propertyName, flagName, boolValue(propertyValue.Default), description)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "cmd.Flags().StringArrayVar(&%s, \"%s\", nil, \"%s\")\n", propertyName, flagName, description)
				} else {
					fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyName, propertyValue)
				}
			case "unknown":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"%s\")\n", propertyName, flagName, description)
			default:
				fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyName, propertyValue)
			}
		}
	}
	fmt.Fprintf(s, "return cmd\n")
	fmt.Fprintf(s, "}\n")
	if true { // append lexicon source to generated file
		lexicon.generateSourceComment(s)
	}
	if err := writeFormattedFile(filename, s.String()); err != nil {
		log.Errorf("error writing file %s %s", filename, err)
	}
}

func int64Value(v any) int64 {
	switch v := v.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 1
	}
}

func boolValue(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	default:
		return false
	}
}
