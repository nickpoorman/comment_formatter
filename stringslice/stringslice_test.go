package stringslice_test

import (
	"testing"

	"github.com/nickpoorman/comment_formatter/stringslice"
	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	cases := []struct {
		slice []string
		want  [][]string
	}{
		{
			slice: []string{
				"// foo",
				"//",
				"// bar",
				"// ping",
				"//",
			},
			want: [][]string{
				{
					"// foo",
				},
				{
					"// bar",
					"// ping",
				},
				{},
			},
		},
		{
			slice: []string{
				"//",
			},
			want: [][]string{
				{},
				{},
			},
		},
	}

	for _, c := range cases {
		got := stringslice.Split(c.slice, "//")
		if !assert.Equal(t, c.want, got) {
			t.Errorf("stringslice.Split(%q, %q) == %q, expected %q", c.slice, "//", got, c.want)
			t.Errorf("\ngot= %#v\nwant=%#v", got, c.want)
		}
	}
}

func TestSliceWithTrim(t *testing.T) {
	cases := []struct {
		slice []string
		want  [][]string
	}{
		{
			slice: []string{
				"// foo",
				"//",
				"// bar",
				"// ping",
				"//",
			},
			want: [][]string{
				{
					"// foo",
				},
				{
					"// bar",
					"// ping",
				},
			},
		},
		{
			slice: []string{
				"//",
			},
			want: [][]string{},
		},
	}

	for _, c := range cases {
		got := stringslice.Split(c.slice, "//")
		got = stringslice.Trim(got)
		if !assert.Equal(t, c.want, got) {
			t.Errorf("stringslice.Split(%q, %q) == %q, expected %q", c.slice, "//", got, c.want)
			t.Errorf("\ngot= %#v\nwant=%#v", got, c.want)
		}
	}
}

func TestFlatten(t *testing.T) {
	cases := []struct {
		slices [][]string
		want   []string
	}{
		{
			slices: [][]string{
				{
					"// foo",
					"// bar",
				},
				{
					"// ping",
				},
				{
					"// pong",
				},
			},
			want: []string{
				"// foo",
				"// bar",
				"// ping",
				"// pong",
			},
		},
	}

	for _, c := range cases {
		got := stringslice.Flatten(c.slices)
		if !assert.Equal(t, c.want, got) {
			t.Errorf("stringslice.Flatten(%q) == %q, expected %q", c.slices, got, c.want)
			t.Errorf("\ngot= %#v\nwant=%#v", got, c.want)
		}
	}
}
