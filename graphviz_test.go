package graphviz_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/goccy/go-graphviz"
)

func TestGraphviz_Image(t *testing.T) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer func() {
		graph.Close()
		g.Close()
	}()
	n, err := graph.CreateNode("n")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	m, err := graph.CreateNode("m")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	e, err := graph.CreateEdge("e", n, m)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	e.SetLabel("e")

	t.Run("png", func(t *testing.T) {
		t.Run("Render", func(t *testing.T) {
			var buf bytes.Buffer
			if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
				t.Fatalf("%+v", err)
			}
			if len(buf.Bytes()) != 4602 {
				t.Fatalf("failed to encode png: bytes length is %d", len(buf.Bytes()))
			}
		})
		t.Run("RenderImage", func(t *testing.T) {
			image, err := g.RenderImage(graph)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			bounds := image.Bounds()
			if bounds.Max.X != 83 {
				t.Fatal("failed to get image")
			}
			if bounds.Max.Y != 177 {
				t.Fatal("failed to get image")
			}
		})
	})
	t.Run("jpg", func(t *testing.T) {
		t.Run("Render", func(t *testing.T) {
			var buf bytes.Buffer
			if err := g.Render(graph, graphviz.JPG, &buf); err != nil {
				t.Fatalf("%+v", err)
			}
			if len(buf.Bytes()) != 3296 {
				t.Fatalf("failed to encode jpg: bytes length is %d", len(buf.Bytes()))
			}
		})
		t.Run("RenderImage", func(t *testing.T) {
			image, err := g.RenderImage(graph)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			bounds := image.Bounds()
			if bounds.Max.X != 83 {
				t.Fatal("failed to get image")
			}
			if bounds.Max.Y != 177 {
				t.Fatal("failed to get image")
			}
		})
	})
}

func TestParseBytes(t *testing.T) {
	type test struct {
		input          string
		expected_valid bool
	}

	tests := []test{
		{input: "graph test { a -- b }", expected_valid: true},
		{input: "graph test { a -- b", expected_valid: false},
		{input: "graph test { a -- b }", expected_valid: true},
		{input: "graph test { a -- }", expected_valid: false},
		{input: "graph test { a -- c }", expected_valid: true},
		{input: "graph test { a - b }", expected_valid: false},
		{input: "graph test { d -- e }", expected_valid: true},
	}

	for i, test := range tests {
		_, err := graphviz.ParseBytes([]byte(test.input))
		actual_valid := err == nil
		if actual_valid != test.expected_valid {
			t.Errorf("Test %d of TestParseBytes failed. Parsing error: %+v", i+1, err)
		}
	}
}

func TestParseFile(t *testing.T) {
	type test struct {
		input          string
		expected_valid bool
	}

	tests := []test{
		{input: "graph test { a -- b }", expected_valid: true},
		{input: "graph test { a -- b", expected_valid: false},
		{input: "graph test { a -- b }", expected_valid: true},
		{input: "graph test { a -- }", expected_valid: false},
		{input: "graph test { a -- c }", expected_valid: true},
		{input: "graph test { a - b }", expected_valid: false},
		{input: "graph test { d -- e }", expected_valid: true},
	}

	createTempFile := func(t *testing.T, content string) *os.File {
		file, err := ioutil.TempFile("", "*")
		if err != nil {
			t.Fatalf("There was an error creating a temporary file. Error: %+v", err)
			return nil
		}
		_, err = file.WriteString(content)
		if err != nil {
			t.Fatalf("There was an error writing '%s' to a temporary file. Error: %+v", content, err)
			return nil
		}
		return file
	}

	for i, test := range tests {
		tmpfile := createTempFile(t, test.input)
		defer os.Remove(tmpfile.Name())

		_, err := graphviz.ParseFile(tmpfile.Name())
		actual_valid := err == nil
		if actual_valid != test.expected_valid {
			t.Errorf("Test %d of TestParseFile failed. Parsing error: %+v", i+1, err)
		}
	}
}
