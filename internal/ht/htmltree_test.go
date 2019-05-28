package htmltree

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	type items struct {
		e   *ElementTree //what to render
		exp string       // expected result
	}
	table := []items{
		{Html(""), "<html></html>"},
		{P(`class=myclass`), "<p class=myclass></p>"},
		{P(`data-foo="foo text"`), `<p data-foo="foo text"></p>`},
	}
	for _, test := range table {
		var b bytes.Buffer
		test.e.Render(&b, -1)
		r := b.String()
		if r != test.exp {
			t.Errorf("Expected %s, got %s", test.exp, r)
		}
	}
}

func TestRenderError(t *testing.T) {
	var b bytes.Buffer
	var err error
	br := Br("")
	br.C = append(br.C, SC("junk"))
	err = br.Render(&b, -1)
	if err.Error() != "br : empty tag may not have content" {
		t.Errorf("%v", err)
	}
	// now test a nested error
	d := Div("", P("", br))
	err = d.Render(&b, -1)
	if err.Error() != "div : p : br : empty tag may not have content" {
		t.Errorf("%v", err)
	}

}

func BenchmarkRender(b *testing.B) {
	meta := Meta(`title="Demo"`)
	head := Head("id=2 class=foo", meta)
	body := Body("id=3 class=bar", Div("", SC("hello"), Br(``)))
	html := Html("", head, body)

	for i := 0; i < b.N; i++ {
		var b bytes.Buffer
		html.Render(&b, -1)
	}
}
