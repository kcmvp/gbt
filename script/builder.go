package script

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type cQC struct {
	root   string
	target string
	err    error
}

const (
	lineCoverage    = "line_coverage.data"
	methodCoverage  = "method_coverage.data"
	testOutput      = "test.data"
	testReport      = "test.json"
	secJson         = "security.json"
	staticJson      = "static.json"
	report          = "index.html"
	secScanTool     = "github.com/securego/gosec/v2/cmd/gosec@latest"
	staticCheckTool = "honnef.co/go/tools/cmd/staticcheck@latest"
)

type TestCase struct {
	Package string
	Test    string
	Action  string
	Output  string
	Elapsed float64
}

type TestCoverage struct {
	Coverage float64
	Packages map[string]*Package
}

type File struct {
	Methods []*Method
	Changes []int
}

type Method struct {
	Name     string
	Coverage float64
}

type Package struct {
	//Name     string
	Coverage float64
	Elapsed  float64
	Files    map[string]*File
}

func projectRoot() string {
	_, file, _, ok := runtime.Caller(2)
	if ok {
		p := filepath.Dir(file)
		for {
			if _, err := os.ReadFile(filepath.Join(p, "go.mod")); err == nil {
				return p
			} else {
				p = filepath.Dir(p)
			}
		}
	}
	panic("Can't figure out project root directory")
}

func NewCQC() *cQC {
	cqc := &cQC{
		root: projectRoot(),
	}
	cqc.target = filepath.Join(cqc.root, "target")
	return cqc
}

func (cqc *cQC) ProjectRoot() string {
	return cqc.root
}

func (cqc *cQC) BuildTarget() string {
	return cqc.target
}

func (cqc *cQC) Clean() *cQC {
	fmt.Println("Clean target ......")
	os.RemoveAll(cqc.target)
	return cqc
}

func (cqc *cQC) Error() error {
	return cqc.err
}

func generateTestReport(cqc *cQC) {
	file, err := os.Open(filepath.Join(cqc.target, testOutput))
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open the file %v", filepath.Join(cqc.target, testOutput)))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	testCoverage := &TestCoverage{
		Packages: map[string]*Package{},
	}
	for scanner.Scan() {
		text := scanner.Text()
		c := TestCase{}
		json.Unmarshal([]byte(text), &c)
		pkg, ok := testCoverage.Packages[c.Package]
		if !ok {
			pkg = &Package{}
			testCoverage.Packages[c.Package] = pkg
		}
		if len(c.Test) == 0 {
			if strings.HasPrefix(c.Output, "coverage:") {
				pcts := strings.TrimRight(strings.Fields(c.Output)[1], "%")
				if pct, err := strconv.ParseFloat(pcts, 2); err == nil {
					pkg.Coverage = pct
				}
			}
			if c.Elapsed > 0 {
				pkg.Elapsed = c.Elapsed
			}
		}
	}

	mc, err := os.Open(filepath.Join(cqc.target, methodCoverage))
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open the file %v", filepath.Join(cqc.target, methodCoverage)))
	}
	defer mc.Close()
	scanner = bufio.NewScanner(mc)
	for scanner.Scan() {
		text := scanner.Text()
		items := strings.Fields(text)
		coverage, _ := strconv.ParseFloat(strings.TrimRight(items[2], "%"), 2)
		if strings.EqualFold(items[0], "total:") {
			testCoverage.Coverage = coverage
		} else {
			m := &Method{
				Name:     items[1],
				Coverage: coverage,
			}
			for s, p := range testCoverage.Packages {
				fName := strings.Split(items[0], ":")[0]
				pkgName := fName[:strings.LastIndex(fName, "/")]
				if strings.EqualFold(pkgName, s) {
					if f, ok := p.Files[fName]; ok {
						f.Methods = append(f.Methods, m)
					} else {
						p.Files = map[string]*File{}
						p.Files[fName] = &File{Methods: []*Method{m}}
						testCoverage.Packages[pkgName] = p
					}
				}
			}
		}
	}
	data, _ := json.Marshal(testCoverage)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "\t")
	os.WriteFile(filepath.Join(cqc.target, testReport), prettyJSON.Bytes(), os.ModePerm)
}

// Test run the test with -race, -cover, -fuzz and -bench
func (cqc *cQC) Test(args ...string) *cQC {
	//cqc.validate()
	fmt.Println("Running unit tests ......")
	os.Chdir(cqc.root)
	os.MkdirAll(cqc.target, os.ModePerm)
	params := []string{"test", "-v", "-json", "-coverprofile", filepath.Join(cqc.target, lineCoverage), "./..."}
	if len(args) > 0 {
		params = append(params, args...)
	}
	out, err := exec.Command("go", params...).CombinedOutput()
	if err != nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "\"Action\":\"fail\"") {
				fmt.Println(line)
			}
		}
		cqc.err = err
		return cqc
	}
	os.WriteFile(filepath.Join(cqc.target, testOutput), out, os.ModePerm)
	//  go tool cover -func ./target/coverage.data
	params = []string{"tool", "cover", "-func", filepath.Join(cqc.target, lineCoverage)}
	out, _ = exec.Command("go", params...).CombinedOutput()
	os.WriteFile(filepath.Join(cqc.target, methodCoverage), out, os.ModePerm)
	generateTestReport(cqc)
	return cqc
}

// Build walk from project root dir and run build command for each executable
// and place the executable at ${project_root}/bin; in case there are more than one executable
func (cqc *cQC) Build(files ...string) *cQC {
	targetFiles := files
	if len(targetFiles) == 0 {
		targetFiles = append(targetFiles, "main.go")
	}
	fmt.Println("Building project ......")
	os.MkdirAll(cqc.target, os.ModePerm)
	filepath.Walk(cqc.root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, t := range targetFiles {
			if strings.EqualFold(info.Name(), t) {
				if output, err := exec.Command("go", "build", "-o", cqc.target, path).CombinedOutput(); err != nil {
					cqc.err = err
					fmt.Println(string(output))
				}
			}
		}
		return nil
	})
	return cqc
}

func (cqc *cQC) SecScan() *cQC {
	fmt.Println("Scanning security issues ......")
	_, err := exec.Command("gosec", "-version").CombinedOutput()
	if err != nil {
		exec.Command("go", "install", secScanTool).CombinedOutput()
	}

	os.MkdirAll(cqc.target, os.ModePerm)

	exec.Command("gosec", "-fmt", "json", "-out", filepath.Join(cqc.target, secJson), "./...").CombinedOutput()
	return cqc
}

func (cqc *cQC) StaticScan() *cQC {
	fmt.Println("Analyzing code ......")
	_, err := exec.Command("staticcheck", "-version").CombinedOutput()
	if err != nil {
		exec.Command("go", "install", staticCheckTool).CombinedOutput()
	}

	os.MkdirAll(cqc.target, os.ModePerm)
	out, _ := exec.Command("staticcheck", "-f", "json", "./...").CombinedOutput()
	items := strings.Split(strings.Trim(string(out), "\n"), "\n")
	result := fmt.Sprintf("{\"Total\":%d, \"Issues\": [%s]}", len(items), strings.Join(items, ","))

	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, []byte(result), "", "\t"); err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
	}
	os.WriteFile(filepath.Join(cqc.target, staticJson), prettyJSON.Bytes(), os.ModePerm)
	return cqc
}

func (cqc *cQC) NotCoveredModified() string {
	cqc.Clean().Test()
	return ""
}
