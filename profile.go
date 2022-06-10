package env

import (
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

const configFileName = "application"
const testContext = "test"
const defaultContext = "default"

var profile *Profile
var once sync.Once

type Resource interface {
	Key() string
	Object() interface{}
}

var ActiveProfile = func() *Profile {
	once.Do(func() {
		// call stack is expected less than 30
		pcs := make([]uintptr, 30)
		n := runtime.Callers(0, pcs)
		pcs = pcs[:n]
		frames := runtime.CallersFrames(pcs)
		for {
			frame, more := frames.Next()
			if !more {
				break
			}
			isCallFromTest, _ := regexp.MatchString(".*_test\\.go$", frame.File)
			if isCallFromTest {
				withProfile(testContext)
				break
			}
		}
	})
	return profile
}

func init() {
	profile = &Profile{
		viper.New(),
		afero.NewOsFs(),
		defaultContext,
	}
	profile.SetConfigName(configFileName)
	profile.SetConfigType("yml")
	profile.AddConfigPath(".")
	err := profile.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	log.SetFlags(0)
}

type Profile struct {
	*viper.Viper
	afero.Fs
	name string
}

func withProfile(p string) {
	profile.name = p
	name := fmt.Sprintf("%s-%s", configFileName, p)
	f := searchInPath(".", name)
	file, err := afero.ReadFile(profile.Fs, f)
	if err != nil {
		log.Fatal(err)
	}
	err = profile.MergeConfig(bytes.NewReader(file))
	if err != nil {
		log.Fatal(err)
	}
}

func (p Profile) Name() string {
	return p.name
}

func absPathify(inPath string) string {
	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		inPath = userHomeDir() + inPath[5:]
	}

	inPath = os.ExpandEnv(inPath)

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}
	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func searchInPath(in, name string) (filename string) {
	for _, ext := range viper.SupportedExts {
		if b, _ := exists(profile.Fs, filepath.Join(in, name+"."+ext)); b {
			return filepath.Join(in, name+"."+ext)
		}
	}
	return ""
}

// Check if file Exists
func exists(fs afero.Fs, path string) (bool, error) {
	stat, err := fs.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
