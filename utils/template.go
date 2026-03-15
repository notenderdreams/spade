package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\{(\w+)(?:=([^}]*))?\}`)
var defaultRe = regexp.MustCompile(`^\((\w+)=(.+)\)$`)

func SubstitutePlaceholders(command string, scriptArgs []string, args []string) (string, []string, error) {
	named := map[string]string{}
	positional := []string{}

	for _, a := range args {
		if k, v, ok := strings.Cut(a, "="); ok {
			named[k] = v
		} else {
			positional = append(positional, a)
		}
	}

	defaults := map[string]string{}
	cleanedScriptArgs := make([]string, 0, len(scriptArgs))
	for _, a := range scriptArgs {
		if m := defaultRe.FindStringSubmatch(a); m != nil {
			defaults[m[1]] = m[2]
		} else {
			cleanedScriptArgs = append(cleanedScriptArgs, a)
		}
	}

	allTemplates := append([]string{command}, cleanedScriptArgs...)
	seen := make([]string, 0, 8)
	seenSet := make(map[string]bool, 8)
	for _, t := range allTemplates {
		for _, m := range placeholderRe.FindAllStringSubmatch(t, -1) {
			k := m[1]
			if !seenSet[k] {
				seen = append(seen, k)
				seenSet[k] = true
			}
			if m[2] != "" && defaults[k] == "" {
				defaults[k] = m[2]
			}
		}
	}

	resolved := map[string]string{}
	posIdx := 0
	for _, k := range seen {
		switch {
		case named[k] != "":
			resolved[k] = named[k]
		case posIdx < len(positional):
			resolved[k] = positional[posIdx]
			posIdx++
		case defaults[k] != "":
			resolved[k] = defaults[k]
		default:
			return "", nil, fmt.Errorf("Missing value for  placeholder {%s}", k)
		}
	}

	expand := func(t string) string {
		return placeholderRe.ReplaceAllStringFunc(t, func(m string) string {
			key := placeholderRe.FindStringSubmatch(m)[1]
			return resolved[key]
		})
	}

	if len(seen) == 0 {
		return command, append(cleanedScriptArgs, args...), nil
	}

	expandedCmd := expand(command)
	expandedArgs := make([]string, len(cleanedScriptArgs))
	for i, a := range cleanedScriptArgs {
		expandedArgs[i] = expand(a)
	}

	expandedArgs = append(expandedArgs, positional[posIdx:]...)
	return expandedCmd, expandedArgs, nil
}

func RenderArg(s string) string {
	return placeholderRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := placeholderRe.FindStringSubmatch(m)
		name := parts[1]
		def := parts[2]

		out := PlaceholderBraceStyle.Render("{") + PlaceholderNameStyle.Render(name)
		if def != "" {
			out += PlaceholderEqStyle.Render("=") + PlaceholderDefaultStyle.Render(def)
		}
		out += PlaceholderBraceStyle.Render("}")
		return out
	})
}
