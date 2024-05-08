package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		commentPrefix string
		lines         []string
		lineNumber    int
		lineLength    int
		expected      []string
	}{
		{
			commentPrefix: "//",
			lines: []string{
				"// This is a comment block.",
				"// It has multiple lines.",
				"//",
				"// There is an empty line above.",
				"//",
				"// This is another comment block.",
				"// It also has multiple lines.",
				"// ",
			},
			lineNumber: 0,
			lineLength: 40,
			expected: []string{
				"// This is a comment block. It has",
				"// multiple lines.",
				"//",
				"// There is an empty line above.",
				"//",
				"// This is another comment block. It",
				"// also has multiple lines.",
			},
		},
		{
			commentPrefix: "//",
			lines: []string{
				"//",
				"// This is a comment block.",
				"// It has multiple lines.",
				"//",
				"// There is an empty line above.",
				"//",
				"// This is another comment block.",
				"// It also has multiple lines.",
			},
			lineNumber: 0,
			lineLength: 40,
			expected: []string{
				"// This is a comment block. It has",
				"// multiple lines.",
				"//",
				"// There is an empty line above.",
				"//",
				"// This is another comment block. It",
				"// also has multiple lines.",
			},
		},
		{
			commentPrefix: "//",
			lines:         []string{},
			lineNumber:    0,
			lineLength:    40,
			expected:      []string{},
		},
		{
			commentPrefix: "//",
			lines: []string{
				"//",
			},
			lineNumber: 0,
			lineLength: 40,
			expected: []string{
				"//",
			},
		},
	}

	for _, c := range cases {
		actual := format(c.commentPrefix, c.lines, c.lineNumber, c.lineLength)
		if len(actual) != len(c.expected) {
			t.Errorf("format(%q, %q, %d, %d) == %q, expected %q", c.commentPrefix, c.lines, c.lineNumber, c.lineLength, actual, c.expected)
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("format(%q, %q, %d, %d) == %q, expected %q", c.commentPrefix, c.lines, c.lineNumber, c.lineLength, actual, c.expected)
			}
		}
	}
}

func TestJoinSubBlocksWithEmptyCommentLines(t *testing.T) {
	block := [][]string{
		{
			"// Line 1",
			"// Line 2",
		},
		{
			"// Line 3",
		},
	}

	expected := [][]string{
		{
			"// Line 1",
			"// Line 2",
		},
		{
			"//",
		},
		{
			"// Line 3",
		},
	}

	actual := joinSubBlocksWithEmptyCommentLines("//", block)
	if len(actual) != len(expected) {
		t.Errorf("joinSubBlocksWithEmptyCommentLines(%q, %q) == %q, expected %q", "//", block, actual, expected)
	}
	if !assert.Equal(t, expected, actual) {
		t.Errorf("joinSubBlocksWithEmptyCommentLines(%q, %q) == %q, expected %q", "//", block, actual, expected)
		t.Errorf("\ngot= %#v\nwant=%#v", actual, expected)
	}
}

func TestSplitCommentBlock(t *testing.T) {
	cases := []struct {
		lines    []string
		expected [][]string
	}{
		{
			lines: []string{
				"//",
				"// Line 1",
				"//",
				"// Line 2",
				"// Line 3",
			},
			expected: [][]string{
				{
					"// Line 1",
				},
				{
					"// Line 2",
					"// Line 3",
				},
			},
		},
	}

	for _, c := range cases {
		got := splitCommentBlock("//", c.lines)
		if len(got) != len(c.expected) {
			t.Errorf("splitCommentBlock(%q, %q) == %q, expected %q", "//", c.lines, got, c.expected)
		}
		if !assert.Equal(t, c.expected, got) {
			t.Errorf("splitCommentBlock(%q, %q) == %q, expected %q", "//", c.lines, got, c.expected)
			t.Errorf("\ngot= %#v\nwant=%#v", got, c.expected)
		}
	}
}
