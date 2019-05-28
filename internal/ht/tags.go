package htmltree

// Wrappers for html element tags.
//
// Functions are grouped in the categories given at
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element and
// are alphabetical within groups.
//
// Conventions:
//     Functions are named by tag with initial caps, e.g. Html()
//
//     The signature for non-empty tags is Tagname(a string, c ...Content) *ElementTree
//     The signature for empty tags is Tagname(a string) *ElementTree
//
//     Empty refers to elements that enclose no content and need no closing tag.
//
//     <style> is the only exception. It's signature is Style(**content). More
//     details and explanation in the doc string.
//
// Obsolete and Deprecated Elements:
// No pull requests will be accepted for
// acronym, applet, basefont, big, blink, center, command, content,
// dir, element, font, frame, frameset, isindex, keygen, listing,
// marquee, multicol, nextid, noembed, plaintext, shadow, spacer,
// strike, tt, xmp .

// Main Root

func Html(a string, c ...Content) *ElementTree {
	return &ElementTree{"html", a, c, false}

}

// Document Metadata
// TODO: base

func Head(a string, c ...Content) *ElementTree {
	return &ElementTree{"head", a, c, false}
}

func Body(a string, c ...Content) *ElementTree {
	return &ElementTree{"body", a, c, false}
}

func Link(a string) *ElementTree {
	return &ElementTree{"link", a, []Content{}, true}
}

func Meta(a string) *ElementTree {
	return &ElementTree{"meta", a, []Content{}, true}
}

func Title(a string, c ...Content) *ElementTree {
	return &ElementTree{"title", a, c, false}
}

// Style is a special case in the sense that the only
// valid content is one or more strings of CSS. At this time
// there's no check to complain about other content.
func Style(a string, c ...Content) *ElementTree {
	return &ElementTree{"style", a, c, false}
}

// Content Sectioning
// TODO hgroup

func Address(a string, c ...Content) *ElementTree {
	return &ElementTree{"address", a, c, false}
}

func Article(a string, c ...Content) *ElementTree {
	return &ElementTree{"article", a, c, false}
}

func Aside(a string, c ...Content) *ElementTree {
	return &ElementTree{"aside", a, c, false}
}

func Footer(a string, c ...Content) *ElementTree {
	return &ElementTree{"footer", a, c, false}
}

func Header(a string, c ...Content) *ElementTree {
	return &ElementTree{"header", a, c, false}
}

func H1(a string, c ...Content) *ElementTree {
	return &ElementTree{"h1", a, c, false}
}

func H2(a string, c ...Content) *ElementTree {
	return &ElementTree{"h2", a, c, false}
}

func H3(a string, c ...Content) *ElementTree {
	return &ElementTree{"h3", a, c, false}
}

func H4(a string, c ...Content) *ElementTree {
	return &ElementTree{"h4", a, c, false}
}

func H5(a string, c ...Content) *ElementTree {
	return &ElementTree{"h5", a, c, false}
}

func H6(a string, c ...Content) *ElementTree {
	return &ElementTree{"h6", a, c, false}
}

func Nav(a string, c ...Content) *ElementTree {
	return &ElementTree{"nav", a, c, false}
}

func Section(a string, c ...Content) *ElementTree {
	return &ElementTree{"section", a, c, false}

}

// Text Content

func Blockquote(a string, c ...Content) *ElementTree {
	return &ElementTree{"blockquote", a, c, false}
}

func Dd(a string, c ...Content) *ElementTree {
	return &ElementTree{"dd", a, c, false}
}

func Div(a string, c ...Content) *ElementTree {
	return &ElementTree{"div", a, c, false}
}

func Dl(a string, c ...Content) *ElementTree {
	return &ElementTree{"dl", a, c, false}
}

func Dt(a string, c ...Content) *ElementTree {
	return &ElementTree{"dt", a, c, false}
}

func Figcaption(a string, c ...Content) *ElementTree {
	return &ElementTree{"figcaption", a, c, false}
}

func Figure(a string, c ...Content) *ElementTree {
	return &ElementTree{"figure", a, c, false}
}

func Hr(a string) *ElementTree {
	return &ElementTree{"hr", a, []Content{}, true}
}

func Li(a string, c ...Content) *ElementTree {
	return &ElementTree{"li", a, c, false}
}

func Main(a string, c ...Content) *ElementTree {
	return &ElementTree{"main", a, c, false}
}

func Ol(a string, c ...Content) *ElementTree {
	return &ElementTree{"ol", a, c, false}
}

func P(a string, c ...Content) *ElementTree {
	return &ElementTree{"p", a, c, false}
}

func Pre(a string, c ...Content) *ElementTree {
	return &ElementTree{"pre", a, c, false}
}

func Ul(a string, c ...Content) *ElementTree {
	return &ElementTree{"ul", a, c, false}

}

// Inline Text Semantics
// TODO abbr, bdi, bdo, data, dfn, kbd, mark, q, rp, rt, rtc, ruby,
//      time, var, wbr

func A(a string, c ...Content) *ElementTree {
	return &ElementTree{"a", a, c, false}
}

func B(a string, c ...Content) *ElementTree {
	return &ElementTree{"b", a, c, false}
}

func Br(a string) *ElementTree {
	return &ElementTree{"br", a, []Content{}, true}
}

func Cite(a string, c ...Content) *ElementTree {
	return &ElementTree{"cite", a, c, false}
}

func Code(a string, c ...Content) *ElementTree {
	return &ElementTree{"code", a, c, false}
}

func Em(a string, c ...Content) *ElementTree {
	return &ElementTree{"em", a, c, false}
}

func I(a string, c ...Content) *ElementTree {
	return &ElementTree{"i", a, c, false}
}

func S(a string, c ...Content) *ElementTree {
	return &ElementTree{"s", a, c, false}
}

func Samp(a string, c ...Content) *ElementTree {
	return &ElementTree{"samp", a, c, false}
}

func Small(a string, c ...Content) *ElementTree {
	return &ElementTree{"small", a, c, false}
}

func Span(a string, c ...Content) *ElementTree {
	return &ElementTree{"span", a, c, false}
}

func Strong(a string, c ...Content) *ElementTree {
	return &ElementTree{"strong", a, c, false}
}

func Sub(a string, c ...Content) *ElementTree {
	return &ElementTree{"sub", a, c, false}
}

func Sup(a string, c ...Content) *ElementTree {
	return &ElementTree{"sup", a, c, false}
}

func U(a string, c ...Content) *ElementTree {
	return &ElementTree{"u", a, c, false}

}

// Image and Multimedia

func Area(a string) *ElementTree {
	return &ElementTree{"area", a, []Content{}, true}
}

func Audio(a string, c ...Content) *ElementTree {
	return &ElementTree{"audio", a, c, false}
}

func Img(a string) *ElementTree {
	return &ElementTree{"img", a, []Content{}, true}
}

func Map(a string, c ...Content) *ElementTree {
	return &ElementTree{"map", a, c, false}
}

func Track(a string) *ElementTree {
	return &ElementTree{"track", a, []Content{}, true}
}

func Video(a string, c ...Content) *ElementTree {
	return &ElementTree{"video", a, c, false}

}

// Embedded Content

func Embed(a string) *ElementTree {
	return &ElementTree{"embed", a, []Content{}, true}
}

func Object(a string, c ...Content) *ElementTree {
	return &ElementTree{"object", a, c, false}
}

func Param(a string) *ElementTree {
	return &ElementTree{"param", a, []Content{}, true}
}

func Source(a string) *ElementTree {
	return &ElementTree{"source", a, []Content{}, true}

}

// Scripting

func Canvas(a string, c ...Content) *ElementTree {
	return &ElementTree{"canvas", a, c, false}
}

func Noscript(a string, c ...Content) *ElementTree {
	return &ElementTree{"noscript", a, c, false}
}

func Script(a string, c ...Content) *ElementTree {
	return &ElementTree{"script", a, c, false}

}

// Demarcating Edits
// TODO del, ins

// Table Content
// TODO colgroup (maybe. It's poorly supported.)

func Caption(a string, c ...Content) *ElementTree {
	return &ElementTree{"caption", a, c, false}
}

func Col(a string) *ElementTree {
	return &ElementTree{"col", a, []Content{}, true}
}

func Table(a string, c ...Content) *ElementTree {
	return &ElementTree{"table", a, c, false}
}

func Tbody(a string, c ...Content) *ElementTree {
	return &ElementTree{"tbody", a, c, false}
}

func Td(a string, c ...Content) *ElementTree {
	return &ElementTree{"td", a, c, false}
}

func Tfoot(a string, c ...Content) *ElementTree {
	return &ElementTree{"tfoot", a, c, false}
}

func Th(a string, c ...Content) *ElementTree {
	return &ElementTree{"th", a, c, false}
}

func Thead(a string, c ...Content) *ElementTree {
	return &ElementTree{"thead", a, c, false}
}

func Tr(a string, c ...Content) *ElementTree {
	return &ElementTree{"tr", a, c, false}

}

// Forms

func Button(a string, c ...Content) *ElementTree {
	return &ElementTree{"button", a, c, false}
}

func Datalist(a string, c ...Content) *ElementTree {
	return &ElementTree{"datalist", a, c, false}
}

func Fieldset(a string, c ...Content) *ElementTree {
	return &ElementTree{"fieldset", a, c, false}
}

func Form(a string, c ...Content) *ElementTree {
	return &ElementTree{"form", a, c, false}
}

func Input(a string) *ElementTree {
	return &ElementTree{"input", a, []Content{}, true}
}

func Label(a string, c ...Content) *ElementTree {
	return &ElementTree{"label", a, c, false}
}

func Legend(a string, c ...Content) *ElementTree {
	return &ElementTree{"legend", a, c, false}
}

func Meter(a string, c ...Content) *ElementTree {
	return &ElementTree{"meter", a, c, false}
}

func Optgroup(a string, c ...Content) *ElementTree {
	return &ElementTree{"optgroup", a, c, false}
}

func Option(a string, c ...Content) *ElementTree {
	return &ElementTree{"option", a, c, false}
}

func Output(a string, c ...Content) *ElementTree {
	return &ElementTree{"output", a, c, false}
}

func Progress(a string, c ...Content) *ElementTree {
	return &ElementTree{"progress", a, c, false}
}

func Select(a string, c ...Content) *ElementTree {
	return &ElementTree{"select", a, c, false}
}

func Textarea(a string, c ...Content) *ElementTree {
	return &ElementTree{"textarea", a, c, false}

}

// Interactive Elememts (Experimental. Omitted for now.)

// Web Components (Experimental. Omitted for now.)
