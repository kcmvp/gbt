package script

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kcmvp/gbt/gbtc/cmd/common"
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
	Coverage    float64
	Packages    map[string]*Package
	err         error
	maxCoverage float64
	minCoverage float64
	moduleDir   string
	rootDir     string
	scriptsDir  string
	targetDir   string
	ctx         context.Context
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

var buildTarget = "target"

type TestCase struct {
	Package string
	Test    string
	Action  string
	Output  string
	Elapsed float64
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

func moduleDir() string {
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
	panic("Can't figure out module directory")
}

func NewCQC(coverages ...float64) *cQC {
	cqc := &cQC{
		moduleDir:   moduleDir(),
		minCoverage: -1,
		maxCoverage: -1,
		Packages:    map[string]*Package{},
		ctx:         context.Background(),
	}
	if len(coverages) == 1 {
		cqc.minCoverage = coverages[0]
		cqc.maxCoverage = 100
	} else if len(coverages) >= 2 {
		cqc.minCoverage = coverages[0]
		cqc.maxCoverage = coverages[1]
	}
	if cqc.minCoverage > 0 && cqc.minCoverage >= cqc.maxCoverage {
		log.Fatalf("invalid coverage range %f ~ %f", cqc.minCoverage, cqc.maxCoverage)
	}
	cqc.targetDir = filepath.Join(cqc.moduleDir, buildTarget)
	cqc.scriptsDir = filepath.Join(cqc.moduleDir, common.ScriptDir)
	os.MkdirAll(cqc.targetDir, os.ModePerm)
	os.MkdirAll(cqc.scriptsDir, os.ModePerm)
	return cqc
}

func (cqc *cQC) WithFlag(arg, value interface{}) {
	cqc.ctx = context.WithValue(cqc.ctx, arg, value)
}

func (cqc *cQC) FlagValue(arg interface{}) interface{} {
	return cqc.ctx.Value(arg)
}

func (cqc *cQC) ModuleDir() string {
	return cqc.moduleDir
}

func (cqc *cQC) RootDir() (string, error) {
	var err error
	cqc.rootDir, err = common.ProjectRoot(cqc.moduleDir)
	return cqc.rootDir, err
}

func (cqc *cQC) TargetDir() string {
	return cqc.targetDir
}

func (cqc *cQC) Clean() *cQC {
	fmt.Println("clean build directory ......")
	os.RemoveAll(cqc.targetDir)
	return cqc
}

func (cqc *cQC) Error() error {
	return cqc.err
}
func (cqc *cQC) validateCoverage() {
	if cqc.minCoverage < 0 {
		return
	}
	if cqc.Coverage < cqc.minCoverage {
		cqc.err = fmt.Errorf("miss min coverage %f > %f", cqc.minCoverage, cqc.Coverage)
		return
	}
	f, _ := os.OpenFile(filepath.Join(cqc.scriptsDir, common.CoverageData), os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if c, err := strconv.ParseFloat(s, 64); err == nil {
			if cqc.Coverage < cqc.minCoverage {
				cqc.err = fmt.Errorf("coverage %f less than min %f", cqc.Coverage, cqc.minCoverage)
			} else if cqc.Coverage != c {
				// update only when current coverage > recorded coverage
				f.Truncate(0)
				f.WriteString(strconv.FormatFloat(c, 'f', 4, 64))
			}
		} else {
			cqc.err = err
		}
	}

}

func (cqc *cQC) processResult() {
	file, err := os.Open(filepath.Join(cqc.targetDir, testOutput))
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open the file %v", filepath.Join(cqc.targetDir, testOutput)))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		c := TestCase{}
		json.Unmarshal([]byte(text), &c)
		pkg, ok := cqc.Packages[c.Package]
		if !ok {
			pkg = &Package{}
			cqc.Packages[c.Package] = pkg
		}
		if len(c.Test) == 0 {
			if strings.HasPrefix(c.Output, "coverage:") {
				pcts := strings.TrimRight(strings.Fields(c.Output)[1], "%")
				if pct, err := strconv.ParseFloat(pcts, 64); err == nil {
					pkg.Coverage = pct
				}
			}
			if c.Elapsed > 0 {
				pkg.Elapsed = c.Elapsed
			}
		}
	}

	mc, err := os.Open(filepath.Join(cqc.targetDir, methodCoverage))
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open the file %v", filepath.Join(cqc.targetDir, methodCoverage)))
	}
	defer mc.Close()
	scanner = bufio.NewScanner(mc)
	for scanner.Scan() {
		text := scanner.Text()
		items := strings.Fields(text)
		coverage, _ := strconv.ParseFloat(strings.TrimRight(items[2], "%"), 64)
		if strings.EqualFold(items[0], "total:") {
			cqc.Coverage = coverage
		} else {
			m := &Method{
				Name:     items[1],
				Coverage: coverage,
			}
			for s, p := range cqc.Packages {
				fName := strings.Split(items[0], ":")[0]
				pkgName := fName[:strings.LastIndex(fName, "/")]
				if strings.EqualFold(pkgName, s) {
					if f, ok := p.Files[fName]; ok {
						f.Methods = append(f.Methods, m)
					} else {
						p.Files = map[string]*File{}
						p.Files[fName] = &File{Methods: []*Method{m}}
						cqc.Packages[pkgName] = p
					}
				}
			}
		}
	}
	data, _ := json.Marshal(cqc)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "\t")
	os.WriteFile(filepath.Join(cqc.targetDir, testReport), prettyJSON.Bytes(), os.ModePerm)
}

// Test run the test with -race, -cover, -fuzz and -bench
func (cqc *cQC) Test(args ...string) *cQC {
	fmt.Println("run unit test ......")
	os.Chdir(cqc.moduleDir)
	os.MkdirAll(cqc.targetDir, os.ModePerm)
	params := []string{"test", "-v", "-json", "-coverprofile", filepath.Join(cqc.targetDir, lineCoverage), "./..."}
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
	os.WriteFile(filepath.Join(cqc.targetDir, testOutput), out, os.ModePerm)
	//  go tool cover -func ./targetDir/coverage.data
	params = []string{"tool", "cover", "-func", filepath.Join(cqc.targetDir, lineCoverage)}
	out, _ = exec.Command("go", params...).CombinedOutput()
	os.WriteFile(filepath.Join(cqc.targetDir, methodCoverage), out, os.ModePerm)
	cqc.processResult()
	cqc.validateCoverage()
	return cqc
}

// Build walk from module directory and run build command for each executable
// and place the executable at ${project_root}/bin; in case there are more than one executable
func (cqc *cQC) Build(files ...string) *cQC {
	if cqc.FlagValue(common.FlagBuild).(bool) {
		log.Println("ignore build")
		return cqc
	}
	if cqc.err != nil {
		log.Fatalf("Runs into error %v", cqc.err)
	}
	targetFiles := files
	if len(targetFiles) == 0 {
		targetFiles = append(targetFiles, "main.go")
	}
	fmt.Println("build project ......")
	os.MkdirAll(cqc.targetDir, os.ModePerm)
	filepath.Walk(cqc.moduleDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, t := range targetFiles {
			if strings.EqualFold(info.Name(), t) {
				if output, err := exec.Command("go", "build", "-o", cqc.targetDir, path).CombinedOutput(); err != nil {
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
	fmt.Println("scan security issues ......")
	_, err := exec.Command("gosec", "-version").CombinedOutput()
	if err != nil {
		exec.Command("go", "install", secScanTool).CombinedOutput()
	}
	os.MkdirAll(cqc.targetDir, os.ModePerm)
	exec.Command("gosec", "-fmt", "json", "-out", filepath.Join(cqc.targetDir, secJson), "./...").CombinedOutput()
	return cqc
}

func (cqc *cQC) StaticScan() *cQC {
	fmt.Println("analyze code ......")
	_, err := exec.Command("staticcheck", "-version").CombinedOutput()
	if err != nil {
		exec.Command("go", "install", staticCheckTool).CombinedOutput()
	}
	os.MkdirAll(cqc.targetDir, os.ModePerm)
	out, _ := exec.Command("staticcheck", "-f", "json", "./...").CombinedOutput()
	items := strings.Split(strings.Trim(string(out), "\n"), "\n")
	result := fmt.Sprintf("{\"Total\":%d, \"Issues\": [%s]}", len(items), strings.Join(items, ","))
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, []byte(result), "", "\t"); err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
	}
	os.WriteFile(filepath.Join(cqc.targetDir, staticJson), prettyJSON.Bytes(), os.ModePerm)
	return cqc
}

func (cqc *cQC) NotCoveredModified() string {
	cqc.Clean().Test()
	return ""
}
