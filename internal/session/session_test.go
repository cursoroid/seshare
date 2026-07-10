package session

import "testing"

func TestMungedCwd(t *testing.T) {
	cases := map[string]string{
		"/Users/prathameshmudgale/Work":           "-Users-prathameshmudgale-Work",
		"/Users/prathameshmudgale/Work/qualif-ai": "-Users-prathameshmudgale-Work-qualif-ai",
		"/tmp/se.st_dir-x":                        "-tmp-se-st-dir-x",
		"/a/-Users/b":                             "-a--Users-b",
	}
	for in, want := range cases {
		if got := mungedCwd(in); got != want {
			t.Errorf("mungedCwd(%q) = %q, want %q", in, got, want)
		}
	}
}
