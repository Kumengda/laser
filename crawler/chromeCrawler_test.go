package crawler

import (
	"fmt"
	"github.com/Kumengda/easyChromedp/template"
	"testing"
)

func TestAA(t *testing.T) {
	A := []template.FormData{{
		Name:  "A",
		Type:  "GET",
		Value: "BBB",
	},
		{
			Name:  "B",
			Type:  "GET",
			Value: "BBB",
		},
		{
			Name:  "E",
			Type:  "GET",
			Value: "BBB",
		},
		{
			Name:  "D",
			Type:  "GET",
			Value: "BBB",
		}}
	B := []template.FormData{{
		Name:  "B",
		Type:  "GET",
		Value: "BBB",
	},
		{
			Name:  "A",
			Type:  "GET",
			Value: "BBB",
		},
		{
			Name:  "D",
			Type:  "GET",
			Value: "BBB",
		},
		{
			Name:  "C",
			Type:  "GET",
			Value: "BBB",
		}}
	fmt.Println(mapsEqual(A, B))
}
