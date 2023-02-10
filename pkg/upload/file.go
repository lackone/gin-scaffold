package upload

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const (
	TypeImage FileType = iota + 1
	TypeExcel
	TypeTxt
)

func GetFileName(name string) string {
	ext := GetFileExt(name)
	filename := strings.TrimSuffix(name, ext)
	return utils.EncodeMd5(filename) + ext
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	container := global.Engine.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	return configService.GetString("app.upload_save_path")
}

func CheckSavePath(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

func CheckExt(t FileType, name string) bool {
	container := global.Engine.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range configService.GetStringSlice("app.upload_image_allow_ext") {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}
	}
	return false
}

func CheckSize(t FileType, f multipart.File) bool {
	container := global.Engine.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	all, _ := io.ReadAll(f)
	size := len(all)
	switch t {
	case TypeImage:
		if size <= configService.GetInt("app.upload_image_max_size")*1024*1024 {
			return true
		}
	}
	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, mode os.FileMode) error {
	err := os.MkdirAll(dst, mode)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(f *multipart.FileHeader, dst string) error {
	open, err := f.Open()
	if err != nil {
		return err
	}
	defer open.Close()
	create, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer create.Close()
	_, err = io.Copy(create, open)
	return err
}
