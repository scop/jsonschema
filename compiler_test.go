package jsonschema_test

import (
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestMetaschemaResource(t *testing.T) {
	mainSchema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
        "$schema": "http://tmp.com/meta.json",
        "type": "number"
    }`))
	if err != nil {
		t.Fatal(err)
	}

	metaSchema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$vocabulary": {
            "https://json-schema.org/draft/2020-12/vocab/applicator": true,
            "https://json-schema.org/draft/2020-12/vocab/core": true
        },
        "allOf": [
            { "$ref": "https://json-schema.org/draft/2020-12/meta/applicator" },
            { "$ref": "https://json-schema.org/draft/2020-12/meta/core" }
        ]
    }`))
	if err != nil {
		t.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", mainSchema); err != nil {
		t.Fatal(err)
	}
	if err := c.AddResource("http://tmp.com/meta.json", metaSchema); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Compile("schema.json"); err != nil {
		t.Fatal(err)
	}
}

func TestCompileAnchor(t *testing.T) {
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$defs": {
            "x": {
                "$anchor": "a1",
                "type": "number"
            }
        }
    }`))
	if err != nil {
		t.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schema); err != nil {
		t.Fatal(err)
	}
	sch1, err := c.Compile("schema.json#a1")
	if err != nil {
		t.Fatal(err)
	}
	sch2, err := c.Compile("schema.json#/$defs/x")
	if err != nil {
		t.Fatal(err)
	}
	if sch1 != sch2 {
		t.Fatal("schemas did not match")
	}
}

func TestCompileNonStd(t *testing.T) {
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
        "components": {
            "schemas": {
                "foo" : {
                    "$schema": "https://json-schema.org/draft/2020-12/schema",
                    "$defs": {
                        "x": {
                            "$anchor": "a",
                            "type": "number"
                        },
                        "y": {
                            "$id": "http://temp.com/y",
                            "type": "string"
                        }
                    },
                    "oneOf": [
                        { "$ref": "#a" },
                        { "$ref": "http://temp.com/y" }
                    ]
                }
            }
        }
    }`))
	if err != nil {
		t.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schema); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Compile("schema.json#/components/schemas/foo"); err != nil {
		t.Fatal(err)
	}
}