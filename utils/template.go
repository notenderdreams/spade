package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\{(\w+)\}`)

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

	allTemplates := append([]string{command}, scriptArgs...)
	seen := make([]string, 0, 8)
	seenSet := make(map[string]bool, 8)
	for _, t := range allTemplates {
		for _, m := range placeholderRe.FindAllStringSubmatch(t, -1) {
			if k := m[1]; !seenSet[k] {
				seen = append(seen, k)
				seenSet[k] = true
			}
		}
	}

	resolved := map[string]string{}
	posIdx := 0
	for _, k := range seen {
		if v, ok := named[k]; ok {
			resolved[k] = v
		} else if posIdx < len(positional) {
			resolved[k] = positional[posIdx]
			posIdx++
		} else {
			return "", nil, fmt.Errorf("missing value for placeholder {%s}", k)
		}
	}

	expand := func(t string) string {
		return placeholderRe.ReplaceAllStringFunc(t, func(m string) string {
			key := m[1 : len(m)-1]
			return resolved[key]
		})
	}

	if len(seen) == 0 {
		return command, append(scriptArgs, args...), nil
	}

	expandedCmd := expand(command)
	expandedArgs := make([]string, len(scriptArgs))
	for i, a := range scriptArgs {
		expandedArgs[i] = expand(a)
	}

	expandedArgs = append(expandedArgs, positional[posIdx:]...)
	return expandedCmd, expandedArgs, nil
}
