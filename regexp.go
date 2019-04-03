package main

import "regexp"

type pattern struct {
	definitions        *regexp.Regexp
	variableExpression *regexp.Regexp
	expression         *regexp.Regexp
	funcCall           *regexp.Regexp
	funcArgs           *regexp.Regexp
	importSyntx        *regexp.Regexp
}

func newPattern() *pattern {
	re := &pattern{}

	c, _ := regexp.Compile("\\w[\\w]*=.*")
	re.definitions = c

	c, _ = regexp.Compile("[a-zA-Z_\\.]+")
	re.variableExpression = c

	c, _ = regexp.Compile("[0-9][0-9\\.]*")
	re.expression = c

	c, _ = regexp.Compile("[a-zA-Z_]+\\(.*\\)")
	re.funcCall = c

	c, _ = regexp.Compile(`\((.*)\)`)
	re.funcArgs = c

	c, _ = regexp.Compile(`import.*`)
	re.importSyntx = c

	return re
}
