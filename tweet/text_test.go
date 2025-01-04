package tweet

import (
	"testing"
)

func TestGenText(t *testing.T) {
	cases := []struct {
		strings   []string
		answer string
	}{
		{[]string{"a", "b"}, "a b"},
		{[]string{"a"}, "a"},
		{[]string{"a", "a b"}, "a \"a b\""},
	}

	for _, c := range cases {
		if GenText(c.strings) != c.answer {
			t.Errorf("strings: %s, result: %s, answer: %s", c.strings, GenText(c.strings), c.answer)
		}
	}
}

