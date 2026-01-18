package lexica

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/iancoleman/strcase"
)

func (lexicon *Lexicon) generateCheckCommands(root string) {
	for defname, def := range lexicon.Defs {
		switch def.Type {
		case "record":
			lexicon.generateCheckCommandForDef(root, defname, def)
		}
	}
}

func (lexicon *Lexicon) generateCheckCommandForDef(root, defname string, def *Def) {

	var ntomerge int
	{
		parts0 := strings.Split(lexicon.Id, ".")
		if len(parts0) == 4 {
			ntomerge = 2
		} else if len(parts0) == 3 {
			ntomerge = 1
		}
	}

	id := strings.Replace(lexicon.Id, ".", "-", ntomerge) // merge initial segments of the lexicon id
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
	fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
	fmt.Fprintf(s, "cmd := &cobra.Command{\n")
	fmt.Fprintf(s, "Use: \"%s FILENAME\",\n", commandname)
	fmt.Fprintf(s, "Short: common.Truncate(xrpc.%s_Description),\n", handlerName)
	fmt.Fprintf(s, "Long: xrpc.%s_Description,\n", handlerName)
	fmt.Fprintf(s, "Args: cobra.ExactArgs(1),\n")
	fmt.Fprintf(s, "RunE: func(cmd *cobra.Command, args []string) error {\n")
	fmt.Fprintf(s, "untyped, err := common.ReadJSONFile(args[0])\n")
	fmt.Fprintf(s, "m, ok := untyped.(map[string]any)\n")
	fmt.Fprintf(s, "if ok {m[\"createdAt\"] = common.Now()}\n")
	fmt.Fprintf(s, "if err != nil {return err}\n")
	fmt.Fprintf(s, "b, err := json.Marshal(untyped)\n")
	fmt.Fprintf(s, "if err != nil {return err}\n")
	fmt.Fprintf(s, "var record xrpc.%s\n", handlerName)
	fmt.Fprintf(s, "err = json.Unmarshal(b, &record)\n")
	fmt.Fprintf(s, "b, err = json.MarshalIndent(record, \"\", \"  \")\n")
	fmt.Fprintf(s, "if err != nil {return err}\n")
	fmt.Fprintf(s, "fmt.Fprintf(cmd.OutOrStdout(), \"%%s\\n\", string(b))\n")
	fmt.Fprintf(s, "return nil\n")
	fmt.Fprintf(s, "},\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "return cmd\n")
	fmt.Fprintf(s, "}\n")

	if true { // append lexicon source to generated file
		lexicon.generateSourceComment(s)
	}
	if err := writeFormattedFile(filename, s.String()); err != nil {
		log.Errorf("error writing file %s %s", filename, err)
	}
	log.Debugf("generating %s %s %s %s", filename, commandname, handlerName, packagename)
}
