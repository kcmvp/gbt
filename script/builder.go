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
	"strings"
)

type cQC struct {
	root   string
	target string
}

const (
	//target          = "target"
	coverage        = "coverage.data"
	testData        = "test.data"
	testJson        = "test.json"
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

func InstallDependencies() {
	//@todo install the missing dependencies
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
func (cqc *cQC) Test(args ...string) *cQC {
	//cqc.validate()
	fmt.Println("Running unit tests ......")
	os.Chdir(cqc.root)
	os.MkdirAll(cqc.target, os.ModePerm)
	params := []string{"test", "-v", "-json", "-coverprofile", filepath.Join(cqc.target, coverage), "./..."}
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
		os.Exit(1)
	}
	os.WriteFile(filepath.Join(cqc.target, testData), out, os.ModePerm)
	generateTestReport(filepath.Join(cqc.target, testData))
	processCoverage(filepath.Join(cqc.target, coverage))
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
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, t := range targetFiles {
			if strings.EqualFold(info.Name(), t) {
				if output, err := exec.Command("go", "build", "-o", cqc.target, path).CombinedOutput(); err != nil {
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
	out, err := exec.Command("staticcheck", "-f", "json", "./...").CombinedOutput()
	if err != nil {
		log.Fatalf(string(out))
	}

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

func (cqc *cQC) Cyclomatic() error {

	panic("@todo https://github.com/fzipp/gocyclo")
}
