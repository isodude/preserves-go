package auto

import (
	. "github.com/isodude/preserves-go/lib/preserves"
)

func (d Definitions) Add(k Symbol, v Definition) Definitions {
	d[k] = v
	return d
}

func (d DictionaryEntries) Add(value Value, namedSimplePattern NamedSimplePattern) DictionaryEntries {
	d[value] = namedSimplePattern
	return d
}

func (m Modules) Add(modulePath ModulePath, schema Schema) Modules {
	m[modulePath.ToHash()] = schema
	return m
}
