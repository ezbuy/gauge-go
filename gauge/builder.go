package gauge

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ezbuy/gauge-go/constants"
	"github.com/ezbuy/gauge-go/util"
	"github.com/getgauge/common"
)

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getTempDir() string {
	projectRoot := os.Getenv(common.GaugeProjectRootEnv)
	tempGaugeDir := filepath.Join(projectRoot, "gauge_temp")
	tempGaugeDir += strconv.FormatInt(time.Now().UnixNano(), 10)
	if !exists(tempGaugeDir) {
		os.MkdirAll(tempGaugeDir, 0755)
	}
	return tempGaugeDir
}

// LoadGaugeImpls builds the go project and runs the generated go file,
// so that the gauge specific implementations get scanned
func LoadGaugeImpls() error {
	var b bytes.Buffer
	buff := bufio.NewWriter(&b)

	if err := util.RunCommand(os.Stdout, os.Stdout, constants.CommandGo, "build", "./..."); err != nil {
		buff.Flush()
		return fmt.Errorf("Build failed: %s\n", err.Error())
	}

	// get list of all packages in the projectRoot
	if err := util.RunCommand(buff, buff, constants.CommandGo, "list", "./..."); err != nil {
		buff.Flush()
		fmt.Printf("Failed to get the list of all packages: %s\n%s", err.Error(), b.String())
	}

	tempDir := getTempDir()
	defer os.RemoveAll(tempDir)

	gaugeGoMainFile := filepath.Join(tempDir, constants.GaugeTestFileName)
	f, err := os.Create(gaugeGoMainFile)
	if err != nil {
		return fmt.Errorf("Failed to create main file in %s: %s", tempDir, err.Error())
	}

	genGaugeTestFileContents(f, b.String())
	f.Close()
	// Scan gauge methods
	if err := util.RunCommand(os.Stdout, os.Stdout, constants.CommandGo, "test", "-v", gaugeGoMainFile); err != nil {
		return fmt.Errorf("Failed to compile project: %s\nPlease ensure the project is in GOPATH.\n", err.Error())
	}
	return nil
}

func genGaugeTestFileContents(fileWriter io.Writer, importString string) {
	type info struct {
		Imports []string
	}
	var validImports []string
	for _, i := range strings.Fields(importString) {
		if strings.HasPrefix(i, "_") {
			validImports = append(validImports, strings.TrimPrefix(i, "_"))
		} else {
			validImports = append(validImports, i)
		}
	}
	gaugeTestRunnerTpl.Execute(fileWriter, info{Imports: validImports})
}

var gaugeTestRunnerTpl = template.Must(template.New("main").Parse(
	`package gauge_test_bootstrap
import (
	"os"
	"github.com/ezbuy/gauge-go/gauge"
{{range $n, $i := .Imports}}	_ "{{$i}}"
{{end}})
func init() {
	gauge.Run()
	_, w, _ := os.Pipe()
	os.Stderr = w
	os.Stdout = w
}
`))
