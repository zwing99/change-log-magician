package changelog

import "testing"

func TestParseLevel(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"major", "minor", "patch", "unreleased"} {
		if _, err := ParseLevel(value); err != nil {
			t.Fatalf("ParseLevel(%q): %v", value, err)
		}
	}
	if _, err := ParseLevel("banana"); err == nil {
		t.Fatal("expected invalid level error")
	}
}

func TestParseFragmentAllowsReorderedCategories(t *testing.T) {
	t.Parallel()
	fragment, err := ParseFragment("feature.md", "# [minor]\n\n## Fixed\n\n- Correct output\n\n## Added\n\n- New option\n")
	if err != nil {
		t.Fatal(err)
	}
	if fragment.Level != Minor || len(fragment.Entries[Fixed]) != 1 || len(fragment.Entries[Added]) != 1 {
		t.Fatalf("unexpected fragment: %#v", fragment)
	}
}

func TestParseFragmentErrors(t *testing.T) {
	t.Parallel()
	for _, source := range []string{"## Added\n\n- Missing level\n", "# [patch]\n\n## Surprise\n\n- Nope\n", "# [patch]\n"} {
		if _, err := ParseFragment("bad.md", source); err == nil {
			t.Fatalf("expected error for %q", source)
		}
	}
}

func TestRenderFragment(t *testing.T) {
	t.Parallel()
	got := RenderFragment(Fragment{Level: Patch, Entries: map[Category][]string{Fixed: {"Correct regression"}}})
	want := "# [patch]\n\n## Fixed\n\n- Correct regression\n"
	if got != want {
		t.Fatalf("rendered fragment:\n%s\nwant:\n%s", got, want)
	}
}
