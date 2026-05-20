package bump_version

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestChangeVersion(t *testing.T) {
	testCases := []struct {
		in    string
		vtype VersionType
		out   string
	}{
		{"0.4", Major, "1.0"},
		{"0.4.0", Major, "1.0.0"},
		{"1.0", Major, "2.0"},
		{"1", Major, "2"},
		{"1.0.1", Minor, "1.1.0"},
	}
	for _, tt := range testCases {
		v, err := changeVersion(tt.vtype, tt.in)
		if err != nil {
			t.Fatal(err)
		}
		if v.String() != tt.out {
			t.Errorf("changeVersion(%s, %s): got %s, want %s", tt.vtype, tt.in, v.String(), tt.out)
		}
	}
}

func TestBumpInFileInvalidVersionType(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "version.go")
	const src = "package version\n\nconst Version = \"1.2.3\"\n"
	if err := os.WriteFile(filename, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BumpInFile panicked for invalid version type: %v", r)
		}
	}()

	_, err := BumpInFile(VersionType("ptach"), filename)
	if err == nil {
		t.Fatal("BumpInFile returned nil error for invalid version type")
	}
	if !strings.Contains(err.Error(), "invalid version type") {
		t.Fatalf("BumpInFile error = %q, want invalid version type", err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != src {
		t.Fatalf("BumpInFile changed file on error:\ngot:\n%s\nwant:\n%s", data, src)
	}
}

func TestVersionString(t *testing.T) {
	typ := reflect.TypeOf(VERSION)
	if typ.String() != "string" {
		t.Errorf("expected VERSION to be a string, got %#v (type %#v)", VERSION, typ.String())
	}
}

func TestString(t *testing.T) {
	v := Version{0, 0, 1}
	if v.String() != "0.0.1" {
		t.Errorf("wrong version string reported")
	}
	v = Version{1, 0, 1}
	if v.String() != "1.0.1" {
		t.Errorf("wrong version string reported")
	}
	v = Version{1, -1, -1}
	if want := "1"; v.String() != want {
		t.Errorf("wrong version string reported: got %q want %q", v.String(), want)
	}
}

var lessTests = []struct {
	i, j string
	want bool
}{
	{"1", "2", true},
	{"1.1", "2", true},
	{"1.3.7", "2", true},
	{"1.3.7", "0.1", false},
	{"1.3.7", "1.3.8", true},
	{"1.3.7", "1.3.6", false},
	{"1.3.7", "1.3.7", false},
	{"1.3.7", "1", false},
}

func TestLess(t *testing.T) {
	for _, tt := range lessTests {
		i, _ := Parse(tt.i)
		j, _ := Parse(tt.j)
		got := Less(i, j)
		if got != tt.want {
			t.Errorf("Less(%q, %q): got %t, want %t", i, j, got, tt.want)
		}
	}
}
