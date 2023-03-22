package middleware

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type serializer struct {
	echo.DefaultJSONSerializer
}

func (s *serializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}

	switch v := i.(type) {
	case echo.Map, map[string]interface{}:
		return enc.Encode(v)
	case *Body:
		return enc.Encode(v)
	default:
		return enc.Encode(NewDefaultBody().WithData(i))
	}
}

// JSONSerializer override the default JSON serializer
func JSONSerializer() echo.JSONSerializer {
	return &serializer{}
}
