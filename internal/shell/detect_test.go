package shell

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		input string
		want  Shell
	}{
		{"bash", Bash},
		{"bash-5.2", Bash},
		{"zsh", Zsh},
		{"zsh-5.9", Zsh},
		{"fish", Fish},
		{"dash", Dash},
		{"ksh", Ksh},
		{"mksh", Ksh},
		{"sh", Sh},
		{"BASH", Bash},
		{"Fish", Fish},
		{"", Unknown},
		{"python3", Unknown},
	}

	for _, c := range cases {
		got := parse(c.input)
		if got != c.want {
			t.Errorf("parse(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestShellString(t *testing.T) {
	if Bash.String() != "bash" {
		t.Errorf("Bash.String() = %q, want \"bash\"", Bash.String())
	}
	if Fish.String() != "fish" {
		t.Errorf("Fish.String() = %q, want \"fish\"", Fish.String())
	}
	if Unknown.String() != "unknown" {
		t.Errorf("Unknown.String() = %q, want \"unknown\"", Unknown.String())
	}
}
