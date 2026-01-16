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
		"com.atproto.sync.listRepos",
		"com.atproto.admin.getInviteCodes",
		"com.atproto.server.getAccountInviteCodes",
	}
	if !slices.Contains(allow, lexicon.Id) {
		//return
	}
	for defname, def := range lexicon.Defs {
		if def.Type == "query" {
			lexicon.generateLeafCommandForDef(root, defname, def)
		} else if def.Type == "procedure" {
			lexicon.generateLeafCommandForDef(root, defname, def)
		}
	}
}

func (lexicon *Lexicon) generateLeafCommandForDef(root, defname string, def *Def) {
	if defname != "main" {
		log.Errorf("Can't generate leaf command for %s %s", lexicon.Id, defname)
	}

	id := strings.Replace(lexicon.Id, ".", "-", 1) // merge the first two segments of the lexicon id
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
	fmt.Fprintf(s, "import \"github.com/agentio/slink/api\"\n")
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/common\"\n")
	fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
	if def.Type == "query" && def.Parameters != nil {
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyValue := def.Parameters.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "var %s string\n", propertyName)
			case "integer":
				fmt.Fprintf(s, "var %s int\n", propertyName)
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
	}
	fmt.Fprintf(s, "cmd := &cobra.Command{\n")
	fmt.Fprintf(s, "Use: \"%s\",\n", commandname)
	fmt.Fprintf(s, "Args: cobra.NoArgs,\n")
	fmt.Fprintf(s, "Short: api.%s_Description,\n", handlerName)
	fmt.Fprintf(s, "RunE: func(cmd *cobra.Command, args []string) error {\n")
	if def.Type == "query" && def.Parameters != nil {
		fmt.Fprintf(s, "client := common.NewClient()\n")
		fmt.Fprintf(s, "response, err := api.%s(\n", handlerName)
		fmt.Fprintf(s, "cmd.Context(),\n")
		fmt.Fprintf(s, "client,\n")
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyValue := def.Parameters.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "%s,\n", propertyName)
			case "integer":
				fmt.Fprintf(s, "int64(%s),\n", propertyName)
			case "boolean":
				fmt.Fprintf(s, "%s,\n", propertyName)
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "%s,\n", propertyName)
				}
			default:
			}
		}
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "if err != nil {return err}\n")
		fmt.Fprintf(s, "return common.Write(cmd.OutOrStdout(), response)\n")
	} else {
		fmt.Fprintf(s, "return errors.New(\"unimplemented\")")
	}
	fmt.Fprintf(s, "},\n")
	fmt.Fprintf(s, "}\n")
	if def.Type == "query" && def.Parameters != nil {
		for _, propertyName := range sortedPropertyNames(def.Parameters.Properties) {
			propertyValue := def.Parameters.Properties[propertyName]
			switch propertyValue.Type {
			case "string":
				fmt.Fprintf(s, "cmd.Flags().StringVar(&%s, \"%s\", \"\", \"\")\n", propertyName, propertyName)
			case "integer":
				fmt.Fprintf(s, "cmd.Flags().IntVar(&%s, \"%s\", %d, \"\")\n", propertyName, propertyName, int64Value(propertyValue.Default))
			case "boolean":
				fmt.Fprintf(s, "cmd.Flags().BoolVar(&%s, \"%s\", %t, \"\")\n", propertyName, propertyName, boolValue(propertyValue.Default))
			case "array":
				if propertyValue.Items.Type == "string" {
					fmt.Fprintf(s, "cmd.Flags().StringArrayVar(&%s, \"%s\", nil, \"\")\n", propertyName, propertyName)
				} else {
					fmt.Fprintf(s, "// FIXME cmd.Flags().XXXVar(&%s... %+v\n", propertyName, propertyValue)
				}
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
		return -999
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
