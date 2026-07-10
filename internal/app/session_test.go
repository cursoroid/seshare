package app

import "testing"

func TestMungedCwd(t *testing.T) {
	cases := map[string]string{
		"/Users/prathameshmudgale/Work":          "-Users-prathameshmudgale-Work",
		"/Users/prathameshmudgale/Work/qualif-ai": "-Users-prathameshmudgale-Work-qualif-ai",
		// dots and underscores become dashes
		"/tmp/se.st_dir-x": "-tmp-se-st-dir-x",
		// consecutive non-alnum are NOT collapsed: "/-" -> "--"
		"/a/-Users/b": "-a--Users-b",
	}
	for in, want := range cases {
		if got := mungedCwd(in); got != want {
			t.Errorf("mungedCwd(%q) = %q, want %q", in, got, want)
		}
	}
}
