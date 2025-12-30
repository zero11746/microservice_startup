package env

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Config struct {
	// RootDir 应用根目录
	RootDir string `mapstructure:"root_dir"`
	// ConfDir 应用配置文件根目录
	ConfDir string `mapstructure:"conf_dir"`
	// DataDir 应用数据文件根目录
	DataDir string `mapstructure:"data_dir"`
	// LogDir 应用日志文件根目录
	LogDir string `mapstructure:"log_dir"`
}

var (
	once      sync.Once
	envConfig *Config
)

func GetEnvConfig() *Config {
	return envConfig
}

func InitEnvConfig() error {
	once.Do(func() {
		rootPath := AutoDetectAppRootDir()

		envConfig = &Config{
			RootDir: rootPath,
			ConfDir: path.Join(rootPath, "config"),
			DataDir: path.Join(rootPath, "data"),
			LogDir:  path.Join(rootPath, "log"),
		}
	})
	return nil
}

var AutoDetectAppRootDir = autoDetect

func autoDetect() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	names := []string{
		"go.mod",
		filepath.Join("config", "config.yaml"),
		filepath.Join("config", "app.yaml"),
	}
	dir, err1 := findDirMatch(wd, names)
	if err1 == nil {
		return dir
	}
	return wd
}

var errNotFound = errors.New("cannot found")

// findDirMatch 在指定目录下，向其父目录查找对应的文件是否存在
func findDirMatch(baseDir string, fileNames []string) (dir string, err error) {
	currentDir := baseDir
	for i := 0; i < 20; i++ {
		for _, fileName := range fileNames {
			depsPath := filepath.Join(currentDir, fileName)
			if _, err1 := os.Stat(depsPath); !os.IsNotExist(err1) {
				return currentDir, nil
			}
		}

		currentDir = filepath.Dir(currentDir)

		if currentDir == "." {
			break
		}
	}
	return "", errNotFound
}
