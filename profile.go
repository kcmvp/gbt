package profile

import (
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const ConfigFileName = "application"

var profile *Profile

func init() {
	profile = &Profile{
		viper.New(),
		afero.NewOsFs(),
	}
	profile.SetConfigName(ConfigFileName)
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
}

func With(p string) error {
	name := fmt.Sprintf("%s-%s", ConfigFileName, p)
	f := profile.searchInPath(".", name)
	file, err := afero.ReadFile(profile.Fs, f)
	if err != nil {
		return err
	}
	return profile.MergeConfig(bytes.NewReader(file))
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

func (env *Profile) searchInPath(in, name string) (filename string) {
	for _, ext := range viper.SupportedExts {
		if b, _ := exists(env.Fs, filepath.Join(in, name+"."+ext)); b {
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
