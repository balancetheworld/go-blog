package utils

//统一读取环境变量
import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func init() {
	_ = godotenv.Load()
}

// 读取字符串环境变量；如果环境变量不存在，就返回默认值。
func Get(key string, defaultValue ...string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// 读取环境变量并转成整数；如果为空或转换失败，就返回默认值。
func GetAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

//读取环境变量并转成布尔值；如果为空或转换失败，就返回默认值。
func GetAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return result
}
