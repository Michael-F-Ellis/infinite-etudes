// Atempting to re-create my Python htmltree in go.
package htmltree

import (
	"bytes"
	"fmt"
)

// ElementTree is a struct that represents a tree of nested HTML tags.  It is
// recursive via the Content member, C, which is a list of ElementTree structs
// and/or strings.
type ElementTree struct {
	T     string    // tag name
	A     string    // attributes
	C     []Content // content
	empty bool      // set to true for empty tags like <br>
}

// Content provides an interface to the types that may be present in the C
// member of an ElementTree struct. Such types must define a Render method that
// takes a pointer to a bytes.Buffer that will hold the output and an
// indentation level nindent.
type Content interface {
	Render(b *bytes.Buffer, nindent int) error
}

// SC represents the string Content of a leaf node.
type SC string

// SC.Render() writes the string representation of an SC.
func (sc SC) Render(b *bytes.Buffer, nindent int) error {
	b.WriteString(indentation(nindent))
	b.WriteString(string(sc))
	return nil
}

// *ElementTree.Render() generates the HTML defined by an ElementTree.
func (e *ElementTree) Render(b *bytes.Buffer, nindent int) error {
	var err error
	// render the opening tag
	indent := indentation(nindent)
	b.WriteString(indent)
	b.WriteString("<")
	b.WriteString(e.T)
	// render the attributes
	if len(e.A) > 0 {
		b.WriteString(" ")
	}
	b.WriteString(e.A)
	// close the opening tag
	b.WriteString(">")

	// indentation for nested content.
	rindent := nindent
	if nindent >= 0 {
		rindent = nindent + 1
	}
	if e.empty {
		if len(e.C) == 0 {
			return nil
		} else {
			return fmt.Errorf("%s : empty tag may not have content", e.T)
		}
	}
	// otherwise, recursively render the content
	for _, c := range e.C {
		switch c.(type) {
		case SC:
			c.(SC).Render(b, rindent)
		case *ElementTree:
			err = c.(*ElementTree).Render(b, rindent)
			if err != nil {
				return fmt.Errorf("%s : %v", e.T, err)
			}
		default:
			panic("This should be impossible")
		}
	}
	// render the closing tag
	b.WriteString(indent)
	b.WriteString("</")
	b.WriteString(e.T)
	b.WriteString(">")
	return err
}

// indentation returns a string like "\n  " where the number of spaces is n * 2
// if n is 0 or greater. If n is negative, indentation returns an empty string.
// The negative case supports rendering an entire tree without newlines or
// leading spaces.
func indentation(n int) string {
	if n < 0 {
		return "" // no indentation
	}
	s := "\n"
	for i := 0; i < 2*n; i++ {
		s += " "
	}
	return s
}
