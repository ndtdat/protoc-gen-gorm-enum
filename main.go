package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const funcStr = `func (x *%[1]s) Scan(value any) error {
	*x = %[1]s(%[1]s_value[string(value.([]byte))])

	return nil
}

func (x %[1]s) Value() (driver.Value, error) {
	return x.String(), nil
}

func %[1]s_FromString(str string) (%[1]s, error) {
	value, ok := %[1]s_value[str]
	if !ok {
		return %[1]s(0), fmt.Errorf("cannot parse %[1]s from %%s", str)
	}

	return %[1]s(value), nil
}

func %[1]s_MustParseFromString(str string) %[1]s {
	value, ok := %[1]s_value[str]
	if !ok {
		panic(fmt.Errorf("cannot parse %[1]s from %%s", str))
	}

	return %[1]s(value)
}`

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

// generateFile generates a _gorm_enum.pb.go file containing enum gorm definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	hasEnum := false
	if len(file.Enums) != 0 {
		hasEnum = true
	}

	for _, msg := range file.Messages {
		if len(msg.Enums) != 0 {
			hasEnum = true

			break
		}
	}

	if !hasEnum {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_gorm_enum.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-gorm-enum. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P(`import (
	"database/sql/driver"
	"fmt"
)`)
	g.P()

	for _, enum := range file.Enums {
		g.P(fmt.Sprintf(funcStr, enum.GoIdent.GoName))
		g.P()
	}

	for _, msg := range file.Messages {
		for _, enum := range msg.Enums {
			g.P(fmt.Sprintf(funcStr, enum.GoIdent.GoName))
			g.P()
		}
	}

	return g
}
