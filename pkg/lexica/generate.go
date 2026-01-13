package lexica

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"golang.org/x/tools/imports"
)

func (lexica *Lexica) Generate(root string) error {

	os.RemoveAll(root)

	var wg sync.WaitGroup

	for _, lexicon := range lexica.Lexicons {
		wg.Go(func() {
			packagename, filename := names(root, lexicon.Id)
			log.Infof("%s %s", lexicon.Id, filename)
			if packagename != "" && filename != "" {
				generatefile(filename, packagename, &lexicon)
			}
		})
	}

	wg.Wait()
	return nil
}

func names(root, id string) (string, string) {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		log.Errorf("wtf %s", id)
		return "", ""
	}

	d := root + "/" + parts[0] + "_" + parts[1]

	os.MkdirAll(d, 0755)

	filename := d + "/" + parts[2] + "_" + parts[3] + ".go"

	packagename := parts[0] + "_" + parts[1]

	return packagename, filename
}

func generatefile(filename, packagename string, lexicon *Lexicon) error {
	s := "package " + packagename + "\n\n"

	s += "// " + lexicon.Id + "\n\n"

	s += `import "github.com/agentio/slink/pkg/xrpc"` + "\n"

	prefix := codeprefix(lexicon.Id)

	for name, def := range lexicon.Defs {
		log.Infof("%s %s", name, def.Type)

		var defname string
		if name == "main" {
			defname = prefix
		} else {
			defname = prefix + "_" + capitalize(name)
		}

		switch def.Type {
		case "query":
			if def.Output.Encoding == "application/json" {
				// output
				s += "type " + defname + "_Output struct {\n"
				s += renderproperties(lexicon, def.Output.Schema.Properties, def.Output.Schema.Required)
				s += "}\n\n"
				// parameters
				params := ""
				paramsok := false
				if def.Parameters.Type == "params" {
					s += "// " + fmt.Sprintf("%+v\n", def.Parameters)
					params, paramsok = parseParameters(def.Parameters)
					s += "// " + params + "\n"
				}
				// func
				s += "func " + defname + "(ctx context.Context, c xrpc.Client" + params + ") (*" + defname + "_Output" + ", error) {\n"
				s += "  var output " + defname + "_Output" + "\n"
				s += "params := map[string]interface{}{\n"
				if paramsok {
					for parameterName, _ := range def.Parameters.Properties {
						s += `"` + parameterName + `":` + parameterName + ",\n"
					}
				}
				s += "}\n"
				s += `if err := c.Do(ctx, xrpc.Query, "", "` + lexicon.Id + `", params, nil, &output); err != nil {` + "\n"
				s += "return nil, err\n"
				s += "}\n"
				s += "  return &output, nil\n"
				s += "}\n\n"
			} else {
				s += fmt.Sprintf("// FIXME (query) %+v\n", def)
			}

		case "procedure":
			if def.Output.Encoding == "application/json" || def.Input.Encoding == "application/json" {
				// input
				s += "type " + defname + "_Input struct {\n"
				s += renderproperties(lexicon, def.Input.Schema.Properties, def.Input.Schema.Required)
				s += "}\n\n"
				// output
				s += "type " + defname + "_Output struct {\n"
				s += renderproperties(lexicon, def.Output.Schema.Properties, def.Output.Schema.Required)
				s += "}\n\n"
				// func
				s += "// " + def.Description + "\n"
				s += "func " + defname + "(ctx context.Context, c xrpc.Client, input *" + defname + "_Input) (*" + defname + "_Output" + ", error) {\n"
				s += "  var output " + defname + "_Output" + "\n"
				s += `if err := c.Do(ctx, xrpc.Procedure, "application/json", "` + lexicon.Id + `", nil, input, &output); err != nil {` + "\n"
				s += "return nil, err\n"
				s += "}\n"
				s += "  return &output, nil\n"
				s += "}\n\n"
			} else {
				s += fmt.Sprintf("// FIXME (procedure) %+v\n", def)
			}

		case "object":
			s += "type " + defname + " struct {\n"
			s += renderproperties(lexicon, def.Properties, def.Required)
			s += "}\n\n"

		case "string":
			s += "type " + defname + " string\n"
		}
	}

	if true { // append lexicon source to generated file
		filter := func(s string) string {
			return strings.ReplaceAll(s, "*/*", "[ANY]")
		}
		b, _ := json.MarshalIndent(lexicon, "", "  ")
		s += "/*\n"
		s += filter(string(b)) + "\n"
		s += "*/\n"
	}

	formatted, err := imports.Process(filename, []byte(s), nil)
	if err != nil {
		log.Fatalf("failed to run goimports: %v\n%s", err, s)
	}

	return os.WriteFile(filename, []byte(formatted), 0644)
}

func parseParameters(parameters Parameters) (string, bool) {
	var parms []string
	var parameterNames []string
	for parameterName := range parameters.Properties {
		parameterNames = append(parameterNames, parameterName)
	}
	sort.Strings(parameterNames)
	for _, parameterName := range parameterNames {
		parameterValue := parameters.Properties[parameterName]
		declaration := parameterName + " "
		switch parameterValue.Type {
		case "integer":
			declaration += "int64"
		case "string":
			declaration += "string"
		case "boolean":
			declaration += "bool"
		default:
			return "/* FIXME */", false
		}
		parms = append(parms, declaration)
	}
	return ", " + strings.Join(parms, ", "), true
}

func codeprefix(id string) string {
	parts := strings.Split(id, ".")

	if len(parts) != 4 {
		return ""
	}

	return capitalize(parts[2]) + capitalize(parts[3])

}

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func renderproperties(lexicon *Lexicon, properties map[string]Property, required []string) string {
	var s string

	var propnames []string
	for propname, _ := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)

	for _, propname := range propnames {
		property := properties[propname]
		required := slices.Contains(required, propname)
		switch property.Type {
		case "boolean":
			if required {
				s += capitalize(propname) + " bool `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *bool `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "integer":
			if required {
				s += capitalize(propname) + " int64 `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *int64 `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "string":
			if required {
				s += capitalize(propname) + " string `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *string `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "array":
			itemstype := resolveItemsType(lexicon, property.Items)
			if required {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += "// FIXME: skipping optional array\n"
			}
		default:
			s += "// FIXME: " + propname + " " + fmt.Sprintf("required=%t %+v", required, property) + "\n"
		}
	}
	return s
}

func resolveItemsType(lexicon *Lexicon, items Items) string {
	switch items.Type {
	case "ref":
		ref := items.Ref
		if ref[0] == '#' {
			parts := strings.Split(lexicon.Id, ".")
			if len(parts) != 4 {
				return "/* FIXME: i can't parse this " + lexicon.Id + " */ string"
			}
			typename := capitalize(parts[2]) + capitalize(parts[3]) + "_" + capitalize(ref[1:])
			return "*" + typename
		} else {
			parts := strings.Split(ref, "#")
			if len(parts) != 2 {
				return "/* FIXME " + fmt.Sprintf("%+v", ref) + " */ string"
			}
			id := parts[0]
			tag := parts[1]
			idparts := strings.Split(id, ".")
			if len(idparts) != 4 {
				return "/* FIXME " + fmt.Sprintf("%+v", ref) + " */ string"
			}
			name := capitalize(idparts[2]) + capitalize(idparts[3])
			if tag != "main" {
				name += "_" + capitalize(tag)
			}
			// is the ref target in the same package as the lexicon?
			// if not, we need to add the package name prefix

			if !strings.HasPrefix(lexicon.Id, idparts[0]+"."+idparts[1]+".") {
				prefix := idparts[0] + "_" + idparts[1]
				name = prefix + "." + name
			}

			return "*" + name + "/* " + lexicon.Id + " " + ref + " */"
		}
	default:
	}
	return "/* FIXME */ string"
}
