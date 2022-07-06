package script

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type CQC struct {
	root string
}

//var caller = "script/builder.go"

const (
	target   = "target"
	coverage = "coverage.data"
	testData = "test.data"
	testJson = "test.json"
	report   = "index.html"
)

type TestCase struct {
	Package string
	Test    string
	Action  string
	Output  string
	Elapsed float32
}

type Package struct {
	Name      string
	Coverage  float32
	Failed    int
	Elapsed   float32
	UnCovered []string
	Tests     []*TestCase
}

var pkgMap = make(map[string]*Package)

//func setCaller(c string) {
//	caller = c
//}

func NewCQC() (*CQC, error) {
	cqc := &CQC{}
	_, file, _, ok := runtime.Caller(1)
	if ok {
		p := filepath.Dir(file)
		for {
			if _, err := os.ReadFile(filepath.Join(p, "go.mod")); err == nil {
				cqc.root = p
				break
			} else {
				p = filepath.Dir(p)
			}
		}
	}
	if len(cqc.root) > 0 {
		return cqc, nil
	} else {
		return nil, fmt.Errorf("can not get project root directory")
	}
}

func (b CQC) ProjectRoot() string {
	return b.root
}

func (b *CQC) Clean() error {
	fmt.Println("Clean project...")
	os.RemoveAll(filepath.Join(b.root, target))
	return nil
}

func getPkg(pgkName string) *Package {
	if v, o := pkgMap[pgkName]; o {
		return v
	} else {
		v = &Package{
			Name: pgkName,
		}
		pkgMap[pgkName] = v
		return v
	}
}

func generateTestReport(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open the file %v", path))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	var previous TestCase
	for scanner.Scan() {
		text := scanner.Text()
		c := TestCase{}
		json.Unmarshal([]byte(text), &c)
		pkg := getPkg(c.Package)
		if len(c.Test) == 0 {
			if strings.Contains(c.Output, "coverage:") {
				// @todo parse coverage
				pkg.Coverage = 0
			}
			pkg.Elapsed = c.Elapsed
			continue
		}
		pkg.Tests = append(pkg.Tests, &c)
		if strings.EqualFold(c.Action, "fail") {
			pkg.Failed++
		}
		if c.Test == previous.Test && c.Package == previous.Package {
			previous.Output = previous.Output + c.Output
			previous.Action = c.Action
			previous.Elapsed = c.Elapsed
		} else {
			previous = c
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(fmt.Sprintf("failed to read the file %v, %+v", path, err))
	}
}

func processCoverage(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open the file %v", path))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		entries := strings.Split(text, "/")
		l := strings.TrimSpace(entries[len(entries)-1])
		if strings.HasSuffix(l, "0") {
			pkgName := strings.Join(entries[0:len(entries)-1], "/")
			// @todo corner case : no test at all
			if pkg := getPkg(pkgName); pkg != nil {
				pkg.UnCovered = append(pkg.UnCovered, l)
			}
		}
	}
	d, _ := json.MarshalIndent(pkgMap, "", "\t")
	if e := os.WriteFile(filepath.Join(filepath.Dir(path), testJson), d, os.ModePerm); e != nil {
		log.Fatal(fmt.Sprintf("failed to generate coverage report %+v", e))
	}
}

// Test run the test with -race, -cover, -fuzz and -bench
func (b *CQC) Test(args ...string) error {
	fmt.Println("Test project...")
	os.Chdir(b.root)
	buildDir := filepath.Join(b.root, target)
	os.MkdirAll(buildDir, os.ModePerm)
	params := []string{"test", "-v", "-json", "./...", "-coverprofile", filepath.Join(buildDir, coverage)}
	params = append(params, args...)
	out, err := exec.Command("go", params...).CombinedOutput()
	fmt.Println(string(out))
	os.WriteFile(filepath.Join(buildDir, testData), out, os.ModePerm)
	generateTestReport(filepath.Join(buildDir, testData))
	processCoverage(filepath.Join(buildDir, coverage))
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
