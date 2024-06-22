package cryptadapt

// 但凡 crypt.go 能把 newCipherForConfig 变成 public 的，或者 cipher.go 能把 newCipher 变成 public 的。。。

import (
	"errors"
	"github.com/rclone/rclone/backend/crypt"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/config/obscure"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrorFilePathEmpty   = errors.New("file path is empty")
	ErrorOutputPathEmpty = errors.New("output path is empty")
	ErrorPasswordEmpty   = errors.New("password is empty")
)

var (
	FilenameEncryptionList = []string{crypt.NameEncryptionOff.String()} // 其它选项暂不展示给用户
	FilenameEncodingList   = []string{"base64"}                         // 其它选项暂不展示给用户
)

type configGetter struct {
	password string
	salt     string
}

func (cg configGetter) Get(key string) (value string, ok bool) {
	var m = map[string]string{
		"filename_encryption": "off",
		"password":            obscure.MustObscure(cg.password),
		"password2":           obscure.MustObscure(cg.salt),
		"filename_encoding":   "base64",
	}
	value, ok = m[key]
	return value, ok
}

func getCipher(password string, salt string) (*crypt.Cipher, error) {
	config := configmap.New()
	config.AddGetter(&configGetter{password: password, salt: salt}, configmap.PriorityNormal)
	return crypt.NewCipher(config)
}

func DecryptFile(inputPath string, outputDirPath string, password string, salt string) (tipText string, error error) {
	if inputPath == "" {
		return "", ErrorFilePathEmpty
	}
	if outputDirPath == "" {
		return "", ErrorOutputPathEmpty
	}
	if password == "" {
		return "密码为空", ErrorPasswordEmpty
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return "打开文件失败", err
	}

	cipher, err := getCipher(password, salt)
	if err != nil {
		return "内部错误", err
	}

	decrypted, err := cipher.DecryptData(file)
	if err != nil {
		return "解密失败", err
	}

	var _, inputFileName = filepath.Split(inputPath)
	suffixArray := strings.Split(inputFileName, ".")
	var outputName string
	if len(suffixArray) > 1 && suffixArray[len(suffixArray)-1] == "bin" {
		outputName = strings.Join(suffixArray[0:len(suffixArray)-1], ".")
	} else {
		outputName = strings.Join(suffixArray, ".")
	}

	newFile, err := os.Create(filepath.Join(outputDirPath, outputName))
	if err != nil {
		return "输出文件失败", err
	}

	_, err = io.Copy(newFile, decrypted)
	if err != nil {
		return "输出文件失败", err
	}

	err = newFile.Close()
	if err != nil {
		return "内部错误", err
	}

	return "解密成功", nil
}

func EncryptFile(inputPath string, outputDirPath string, password string, salt string) (string, error) {
	if inputPath == "" {
		return "", ErrorFilePathEmpty
	}
	if outputDirPath == "" {
		return "", ErrorOutputPathEmpty
	}
	if password == "" {
		return "密码为空", ErrorPasswordEmpty
	}

	// 先假设是文件 而不是文件夹
	file, err := os.Open(inputPath)
	if err != nil {
		return "打开文件失败", err
	}

	cipher, err := getCipher(password, salt)
	if err != nil {
		return "内部错误", err
	}

	encrypted, err := cipher.EncryptData(file)
	if err != nil {
		return "加密失败", err
	}

	var _, inputFileName = filepath.Split(inputPath)
	var outputName = inputFileName + ".bin"

	newFile, err := os.Create(filepath.Join(outputDirPath, outputName))
	if err != nil {
		return "输出文件失败", err
	}

	_, err = io.Copy(newFile, encrypted)
	if err != nil {
		return "输出文件失败", err
	}

	err = newFile.Close()
	if err != nil {
		return "内部错误", err
	}

	return "加密成功", nil
}
