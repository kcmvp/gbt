package script

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type CQC struct {
	root string
}

var caller = "script/builder.go"

const (
	target       = "target"
	coverage     = "test_coverage.out"
	testData     = "test.json"
	coverageHtml = "index.html"
)

func setCaller(c string) {
	caller = c
}

func NewCQC() (*CQC, error) {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		fPath := filepath.FromSlash(caller)
		if !strings.HasSuffix(file, fPath) {
			return nil, fmt.Errorf("this method can be called only from ${root}/script/builder.go")
		}
		b := &CQC{
			root: strings.ReplaceAll(file, fPath, ""),
		}
		p := filepath.Dir(file)
		if _, err := os.ReadFile(filepath.Join(p, "go.mod")); err == nil {
			return b, nil
		} else if _, err = os.ReadFile(filepath.Join(b.root, "go.mod")); err == nil {
			return b, nil
		}
	}
	return nil, fmt.Errorf("this method can be called only from ${root}/script/builder.go")
}

func (b CQC) ProjectRoot() string {
	return b.root
}

func (b *CQC) Clean() error {
	fmt.Println("Clean project...")
	os.RemoveAll(filepath.Join(b.root, target))
	return nil
}

// Test run the test with -race, -cover, -fuzz and -bench
func (b *CQC) Test() error {
	fmt.Println("Test project...")
	os.Chdir(b.root)
	buildDir := filepath.Join(b.root, target)
	os.MkdirAll(buildDir, os.ModePerm)
	out, err := exec.Command("go", "test", "-v", "-json", "./...", "-coverprofile", filepath.Join(buildDir, coverage)).CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	os.WriteFile(filepath.Join(buildDir, testData), out, os.ModePerm)
	return err
}

// Build walk from project root dir and run build command for each executable
// and place the executable at ${project_root}/bin; in case there are more than one executable
func (b *CQC) Build() error {
	fmt.Println("Building project ...")
	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
	return cmd.Run()
}

func (b *CQC) SecScan() error {
	//@todo gosec https://opensource.com/article/20/9/gosec
	return nil
}

func (w *CQC) StaticCheck() error {
	panic("@todo https://staticcheck.io/docs/getting-started/")
}

func (b *CQC) Cyclomatic() error {

	panic("@todo https://github.com/fzipp/gocyclo")
}
