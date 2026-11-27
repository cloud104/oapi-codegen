package response

import (
	"cmp"
	"go/ast"
	"go/parser"
	"slices"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type statusCodes []string

func (codes statusCodes) String() string {
	switch len(codes) {
	case 0:
		return ""
	case 1:
		return codes[0]
	default:
		return strings.Join(codes[:len(codes)-1], ", ") + " or " + codes[len(codes)-1]
	}
}

type ResponseGroup struct {
	Type        string
	Method      string
	Fields      []string
	StatusCodes statusCodes
}

func GroupResponses(defs []codegen.ResponseTypeDefinition) []*ResponseGroup {
	byType := make(map[string]*ResponseGroup)

	for _, def := range defs {
		typ := def.Schema.TypeDecl()

		if isAnonymousCompositeType(typ) {
			continue
		}

		group := byType[typ]
		if group == nil {
			name := baseTypeName(typ)
			if def.Schema.ArrayType != nil {
				singular := baseTypeName(def.Schema.ArrayType.TypeDecl())
				name = inflection.Plural(singular)
			}

			group = &ResponseGroup{
				Type:   typ,
				Method: name,
			}
			byType[typ] = group
		}

		group.Fields = append(group.Fields, def.TypeName)
		if !slices.Contains(group.StatusCodes, def.ResponseName) {
			group.StatusCodes = append(group.StatusCodes, def.ResponseName)
		}
	}

	groups := make([]*ResponseGroup, 0, len(byType))
	for _, group := range byType {
		slices.Sort(group.Fields)
		slices.Sort(group.StatusCodes)
		groups = append(groups, group)
	}

	// Map iteration order is undefined, so sort the groups to keep generated
	// output deterministic. Response codes define the primary order, with the
	// response type acting as a stable tie-breaker.
	slices.SortFunc(groups, func(a, b *ResponseGroup) int {
		if n := slices.Compare(a.StatusCodes, b.StatusCodes); n != 0 {
			return n
		}

		return cmp.Compare(a.Type, b.Type)
	})

	return groups
}

func isAnonymousCompositeType(typ string) bool {
	expr, err := parser.ParseExpr(typ)
	if err != nil {
		return false
	}

	// Unwrap pointers and slices/arrays.
	for {
		switch t := expr.(type) {
		case *ast.StarExpr:
			expr = t.X
		case *ast.ArrayType:
			expr = t.Elt
		default:
			switch expr.(type) {
			case *ast.StructType, *ast.MapType, *ast.InterfaceType:
				return true
			default:
				return false
			}
		}
	}
}

func baseTypeName(typ string) string {
	expr, err := parser.ParseExpr(typ)
	if err != nil {
		return typ
	}

	for {
		switch t := expr.(type) {
		case *ast.StarExpr:
			expr = t.X
		case *ast.ArrayType:
			expr = t.Elt
		case *ast.ParenExpr:
			expr = t.X
		case *ast.SelectorExpr:
			return t.Sel.Name
		case *ast.Ident:
			return t.Name
		default:
			return typ
		}
	}
}
