package main

import (
	"os"
	"testing"

	"golang.org/x/net/html"
)

func TestVisitLinks(t *testing.T) {
	// Remove o arquivo CSV criado em testes anteriores
	os.Remove("data/links.csv")

	// Teste o link do Github do usuário michelbernardods
	url := "https://github.com/michelbernardods?tab=repositories"
	visitLinks(url)

	// Verifique se o arquivo CSV foi criado com sucesso
	if _, err := os.Stat("data/links.csv"); os.IsNotExist(err) {
		t.Errorf("file not found: data/links.csv")
	}

	// Verifique se o arquivo CSV contém mais de 0 bytes
	if fileInfo, _ := os.Stat("data/links.csv"); fileInfo.Size() <= 0 {
		t.Errorf("file is empty: data/links.csv")
	}
}

func TestExtractLinks(t *testing.T) {
	// Cria um nó html do tipo ElementNode, com um atributo href válido
	aTag := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{{Key: "href", Val: "https://www.google.com"}},
	}
	// Cria um nó html do tipo ElementNode, sem um atributo href
	bTag := &html.Node{
		Type: html.ElementNode,
		Data: "b",
		Attr: []html.Attribute{{Key: "class", Val: "btn"}},
	}
	// Cria um nó html do tipo ElementNode, com um atributo href inválido
	cTag := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{{Key: "href", Val: "javascript:void(0)"}},
	}

	// Cria um nó html do tipo TextNode
	eTag := &html.Node{
		Type: html.TextNode,
		Data: "Lorem ipsum dolor sit amet",
	}

	testCases := []struct {
		desc     string
		node     *html.Node
		expected bool
	}{
		{
			desc:     "Tag 'a' with valid href",
			node:     aTag,
			expected: true,
		},
		{
			desc:     "Tag 'b' no href attribute",
			node:     bTag,
			expected: false,
		},
		{
			desc:     "Tag 'a' with invalid href",
			node:     cTag,
			expected: false,
		},
		{
			desc:     "Text tag",
			node:     eTag,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			links = []string{} // Limpa a slice de links antes de cada teste
			extractLinks(tc.node)

			if tc.expected && len(links) != 1 {
				t.Errorf("expected 1 link, got %d", len(links))
			}
			if !tc.expected && len(links) != 0 {
				t.Errorf("expected 0 links, got %d", len(links))
			}
		})
	}
}
