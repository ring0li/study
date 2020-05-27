package common

import (
	"fmt"
	"github.com/syyongx/php2go"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//写入app pid 到文件，用于重启，每个app main.go 必须写一个
func WriteAppPidToFile(appName string) {
	pid_path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		ioutil.WriteFile(pid_path+"/"+appName+".pid", []byte(strconv.Itoa(syscall.Getpid())), 0644)
	}
}

//Date 日期格式化
// t Int 只接收32位int 因为一般情况下从数据库取出来都是32位int
// format 如果不传则为默认
// return 格式化时间， 默认为-
func Date(t int, format ...interface{}) string {
	t64 := int64(t)
	tm := time.Unix(t64, 0)
	fformat := "2006-01-02 15:04:05"
	if len(format) > 1 {
		fformat = format[0].(string)
	}
	return tm.Format(fformat)
}

func GenerateMac(data map[string]interface{}, secret string) (mac string) {
	keys := make([]string, 0, len(data))
	for k, _ := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var str string
	for _, key := range keys {
		str += key + data[key].(string)
	}

	// config.Logger.Info("generateMac", zap.String("str", str), zap.String("secret", secret))
	myMac := php2go.Md5(secret + str + secret)
	return myMac
}

//Md5 生成字符串的md5
func Md5(str string) string {
	if len(str) == 0 {
		return php2go.Md5(php2go.Uniqid(string(time.Now().UnixNano())))
	}
	return php2go.Md5(str)
}

// left is null; return right
func CheckStringZero(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// left is null; return right
func GetNotNullString(left, right string) string {
	left = strings.TrimSpace(left)
	if len(left) != 0 {
		return left
	} else {
		return right
	}
}

func GetNotNullPwd(left, right string) string {
	left = strings.TrimSpace(left)
	if len(left) != 0 {
		return Md5(left)
	} else {
		return right
	}
}

//limit [:word:]
func LimitWord(str string) bool {
	b, _ := regexp.Match(`[[:word:]]+`, []byte(str))
	return b
}

//
func GetNewBranchName(projectName, name string) string {
	_uuid := php2go.Uniqid("")
	return fmt.Sprintf("%s_%s%s_%s", time.Now().Format("20060102_150405"), name, strings.Replace(projectName, "/", "_", -1), _uuid[8:])
}

//
func IsExistDir(path string) error {
	if f, err := os.Stat(path); err != nil || f.IsDir() {
		return err
	}
	return nil
}

//
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//
func Command(args ...string) ([]byte, error) {
	//
	_args := args
	// _args := append([]string{"-c"}, args...)
	//_args = append(_args, []string{"||", "exit 1"}...)
	//
	cmd := exec.Command("/bin/bash", _args...)
	//
	fmt.Println(">", _args)
	//
	return cmd.CombinedOutput()
}
