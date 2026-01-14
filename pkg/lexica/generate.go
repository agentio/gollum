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
			log.Debugf("%s %s", lexicon.Id, filename)
			if packagename != "" && filename != "" {
				generatefile(filename, packagename, lexicon)
			}
		})
	}
	wg.Wait()
	return nil
}

func names(root, id string) (string, string) {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		log.Warnf("skipping three-segment name %s", id)
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

	s += `import "github.com/agentio/slink/gen/com_atproto"` + "\n"
	s += `import "github.com/agentio/slink/gen/app_bsky"` + "\n"

	prefix := codeprefix(lexicon.Id)

	for name, def := range lexicon.Defs {

		var defname string
		if name == "main" {
			defname = prefix
		} else {
			defname = prefix + "_" + capitalize(name)
		}

		switch def.Type {
		case "query":
			if def.Output != nil && def.Output.Encoding == "application/json" {
				// output
				s += "type " + defname + "_Output struct {\n"
				s += renderproperties(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
				s += "}\n\n"
				s += renderunions(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
				// parameters
				params := ""
				paramsok := false
				if def.Parameters != nil && def.Parameters.Type == "params" {
					s += "// " + fmt.Sprintf("%+v\n", def.Parameters)
					params, paramsok = parseParameters(def.Parameters)
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
			if def.Output != nil && def.Output.Encoding == "application/json" &&
				def.Input != nil && def.Input.Encoding == "application/json" {
				// input
				s += "type " + defname + "_Input struct {\n"
				s += renderproperties(lexicon, defname+"_Input", def.Input.Schema.Properties, def.Input.Schema.Required)
				s += "}\n\n"
				s += renderunions(lexicon, defname+"_Input", def.Input.Schema.Properties, def.Input.Schema.Required)
				// output
				s += "type " + defname + "_Output struct {\n"
				s += renderproperties(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
				s += "}\n\n"
				s += renderunions(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
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
			s += renderproperties(lexicon, defname, def.Properties, def.Required)
			s += "}\n\n"
			s += renderunions(lexicon, defname, def.Properties, def.Required)

		case "string":
			s += "type " + defname + " string\n"

		case "record":
			s += "type " + defname + " struct {\n"
			s += renderproperties(lexicon, defname, def.Properties, def.Required)
			s += "}\n\n"
			s += renderunions(lexicon, defname, def.Properties, def.Required)

		case "array":
			s += "type " + defname + "_Elem struct {\n"
			s += "}\n\n"

		case "token":
			s += "// " + def.Description + "\n"
			s += "const " + defname + " string = " + `"` + name + `"` + "\n\n"

		default:
			log.Warnf("skipping %s.%s (type %s)", lexicon.Id, name, def.Type)
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

func parseParameters(parameters *Parameters) (string, bool) {
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
		case "array":
			if parameterValue.Items.Type == "string" {
				declaration += "[]string"
			} else {
				return "/* FIXME */", false
			}
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

func renderproperties(lexicon *Lexicon, defname string, properties map[string]Property, required []string) string {
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
				s += capitalize(propname) + " int64 `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *int64 `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "string":
			if required {
				s += capitalize(propname) + " string `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *string `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "array":
			itemstype := resolveItemsType(lexicon, property.Items)
			if required {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "ref":
			reftype := resolveRefType(lexicon, property.Ref)
			s += capitalize(propname) + reftype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
		case "unknown":
			if required {
				s += capitalize(propname) + " interface{} `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *interface{} `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "blob":
			if required {
				s += capitalize(propname) + " []byte `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *[]byte `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "union":
			uniontype := resolveUnionType(lexicon, defname, propname)
			s += capitalize(propname) + " " + uniontype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
		default:
			s += "// FIXME: " + propname + " " + fmt.Sprintf("required=%t %+v", required, property) + "\n"
		}
	}
	return s
}

func renderunions(lexicon *Lexicon, defname string, properties map[string]Property, required []string) string {
	var s string

	var propnames []string
	for propname, _ := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)
	for _, propname := range propnames {
		property := properties[propname]
		switch property.Type {
		case "union":
			uniontype := resolveUnionType(lexicon, defname, propname)
			s += "type " + uniontype + " struct {\n"
			for _, ref := range property.Refs {
				s += "// " + ref + "\n"
			}

			s += "}\n\n"
		}
	}
	return s
}

func resolveUnionType(lexicon *Lexicon, defname, propname string) string {
	return capitalize(defname) + "_" + capitalize(propname) // "string"
}

func resolveItemsType(lexicon *Lexicon, items *Items) string {
	switch items.Type {
	case "string":
		return "string"
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

			return "*" + name
		}
	default:
	}
	return "/* FIXME */ string"
}

func resolveRefType(lexicon *Lexicon, ref string) string {
	if ref[0] == '#' {
		parts := strings.Split(lexicon.Id, ".")
		if len(parts) != 4 {
			return "/* FIXME: i can't parse this " + lexicon.Id + " */ string"
		}
		typename := capitalize(parts[2]) + capitalize(parts[3]) + "_" + capitalize(ref[1:])
		return "*" + typename
	} else {
		parts := strings.Split(ref, "#")

		if len(parts) == 2 || len(parts) == 1 {
			var id string
			var tag string
			if len(parts) == 2 {
				id = parts[0]
				tag = parts[1]
			} else {
				id = parts[0]
				tag = "main"
			}

			var reftype string
			reflexicon := Lookup(id)
			if reflexicon != nil {
				refdef := reflexicon.Lookup(tag)
				if refdef != nil {
					reftype = refdef.Type
				}
			}

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

			if reftype == "array" {
				return "[]" + name + "_Elem"
			}

			return "*" + name
		} else {
			return "/* FIXME ref " + fmt.Sprintf("%+v", ref) + " */ string"
		}
	}
}
