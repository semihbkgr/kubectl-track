package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/semihbkgr/yamldiff/diff"
	"k8s.io/klog"
)

var p printer.Printer

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

// todo
func init() {
	p.LineNumber = false
	p.LineNumberFormat = func(num int) string {
		fn := color.New(color.Bold, color.FgHiWhite).SprintFunc()
		return fn(fmt.Sprintf("%2d | ", num))
	}
	p.Bool = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.Number = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.MapKey = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiCyan),
			Suffix: format(color.Reset),
		}
	}
	p.Anchor = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.Alias = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.String = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiGreen),
			Suffix: format(color.Reset),
		}
	}
}

func YamlRenderString(a map[string]any) string {
	y, err := yaml.Marshal(a)
	if err != nil {
		klog.Warning(err)
		return ""
	}
	tokens := lexer.Tokenize(string(y))
	return p.PrintTokens(tokens)
}

func DiffRenderString(a map[string]any, b map[string]any) string {
	aYaml, err := yaml.Marshal(a)
	if err != nil {
		klog.Warning(err)
		return ""
	}
	bYaml, err := yaml.Marshal(b)
	if err != nil {
		klog.Warning(err)
		return ""
	}

	diffCtx, err := diff.NewDiffContextBytes(aYaml, bYaml, true)
	if err != nil {
		klog.Warning(err)
		return ""
	}

	diffs := diffCtx.Diffs(diff.DefaultDiffOptions)
	return diffs.OutputString(diff.DefaultOutputOptions)
}
