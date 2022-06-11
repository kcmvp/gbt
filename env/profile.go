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
	"sync"
)

const configFileName = "application"
const testProfile = "test"
const defaultProfile = "default"

var profile *Profile
var once sync.Once

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
				withProfile(testProfile)
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
		defaultProfile,
	}
	profile.SetConfigName(configFileName)
	profile.SetConfigType("yml")
	profile.AddConfigPath(".")
	profile.AddConfigPath("../")
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

func (p *Profile) IsTest() bool {
	return p.name == testProfile
}

func withProfile(p string) {
	profile.name = p
	name := fmt.Sprintf("%s-%s", configFileName, p)

	for _, i := range []string{".", "../"} {
		f := searchInPath(i, name)
		if f == "" {
			continue
		}
		file, err := afero.ReadFile(profile.Fs, f)
		if err != nil {
			log.Fatal(err)
		}
		err = profile.MergeConfig(bytes.NewReader(file))
		if err != nil {
			log.Fatal(err)
		} else {
			break
		}
	}
}

func (p Profile) Name() string {
	return p.name
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
