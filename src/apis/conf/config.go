package conf

import (
	"fmt"
	"github.com/joho/godotenv"
	"insightful/src/apis/utils"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
)

// Config struct for config environment
type Config struct {
	DBMysqlUsername     string `env:"DB_MYSQL_USERNAME" default:"root"`
	DBMysqlPassword     string `env:"DB_MYSQL_PASSWORD" default:"root"`
	DBMysqlHost         string `env:"DB_MYSQL_HOST" default:"127.0.0.1"`
	DBMysqlPort         string `env:"DB_MYSQL_PORT" default:"localhost"`
	DBMysqlName         string `env:"DB_MYSQL_NAME" default:"insightfull"`
	DBMysqlMaxIdleConns int    `env:"DB_MYSQL_MAXIDLECONNS" default:"1"`
	DBMysqlMaxOpenConns int    `env:"DB_MYSQL_MAXOPENCONNS" default:"4"`

	LogLevel string `env:"LOG_LEVEL" default:"INFO"`

	MaxWorker int `env:"MAX_WORKER" default:"3"`
	MaxQueue  int `env:"MAX_QUEUE" default:"20"`

	RedisHost     string `env:"REDIS_HOST" default:"127.0.0.1"`
	RedisPort     string `env:"REDIS_PORT" default:"6379"`
	RedisDatabase int    `env:"REDIS_DATABASE" default:"0"`
	RedisPassword string `env:"REDIS_PASSWORD" default:""`
}

// EnvConfig save config from system parameters
var EnvConfig Config

// loadEnvFile load env file from a baseDir, empty mean current working dir
func loadEnvFile(file string) {
	fmt.Println("Loading env from file ", file)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("File", file, "does not exist")
		return
	}
	err := godotenv.Load(file)
	if err != nil {
		panic(err)
	}
}

func loadCwdEnvFile() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	loadEnvFile(filepath.Join(cwd, ".env"))
}

func loadProjectEnvFile() {
	// Load env from project
	_, configFile, _, _ := runtime.Caller(1)
	envFile, err := filepath.Abs(
		filepath.Join(
			filepath.Dir(configFile),
			"..",
			"..",
			".env"))
	if err != nil {
		panic(err)
	}
	loadEnvFile(envFile)
}

func init() {
	loadCwdEnvFile()
	loadProjectEnvFile()

	// parse config
	EnvConfig = Config{}
	err := ParseEnvConfig(&EnvConfig)
	if err != nil {
		fmt.Println(err)
	}
}

// ParseEnvConfig Parse environment to config struct
func ParseEnvConfig(config interface{}) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	var envValue string
	for i := 0; i < t.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)
		envKey := tField.Tag.Get("env")
		if len(envKey) == 0 || envKey == "-" {
			continue
		}
		envValue = os.Getenv(envKey)
		if len(envValue) == 0 {
			envValue = tField.Tag.Get("default")
			if len(envValue) != 0 {
				fmt.Printf("Missng key %s in enviroment. Use \"%s\" instead\n", envKey, envValue)
			} else {
				fmt.Printf("Missng key %s in enviroment\n", envKey)
				if _, ok := tField.Tag.Lookup("require"); ok {
					panic(fmt.Sprintf("ENV \"%s\" is required", envKey))
				}
				continue
			}
		}
		switch vField.Kind() {
		case reflect.Ptr:
			t := vField.Type().String()
			switch t {
			case "*int":
				v, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerInt(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerInt(int(v))))
				break
			case "*int8":
				v, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerInt8(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerInt8(int8(v))))
				break
			case "*int16":
				v, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerInt16(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerInt16(int16(v))))
				break
			case "*int32":
				v, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerInt32(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerInt32(int32(v))))
				break
			case "*int64":
				v, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerInt64(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerInt64(v)))
				break
			case "*uint":
				v, err := strconv.ParseUint(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerUInt(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerUInt(uint(v))))
				break
			case "*uint8":
				v, err := strconv.ParseUint(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerUInt8(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerUInt8(uint8(v))))
				break
			case "*uint16":
				v, err := strconv.ParseUint(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerUInt16(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerUInt16(uint16(v))))
				break
			case "*uint32":
				v, err := strconv.ParseUint(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerUInt32(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerUInt32(uint32(v))))
				break
			case "*uint64":
				v, err := strconv.ParseUint(envValue, 10, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerUInt64(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerUInt64(v)))
				break
			case "*float32":
				v, err := strconv.ParseFloat(envValue, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerFloat32(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerFloat32(float32(v))))
				break
			case "*float64":
				v, err := strconv.ParseFloat(envValue, 64)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerFloat64(0)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerFloat64(v)))
				break
			case "*string":
				vField.Set(reflect.ValueOf(utils.StringToPointer(envValue)))
				break
			case "*bool":
				b, err := strconv.ParseBool(envValue)
				if err != nil {
					fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
					vField.Set(reflect.ValueOf(utils.PointerBoolean(false)))
					break
				}
				vField.Set(reflect.ValueOf(utils.PointerBoolean(b)))
				break
			}
		case reflect.Bool:
			b, err := strconv.ParseBool(envValue)
			if err != nil {
				fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
				vField.SetBool(false)
				break
			}
			vField.SetBool(b)
			break
		case reflect.Int64,
			reflect.Int32,
			reflect.Int16,
			reflect.Int8,
			reflect.Int:
			v, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
				vField.SetInt(0)
				break
			}
			vField.SetInt(v)
			break
		case reflect.Uint64,
			reflect.Uint32,
			reflect.Uint16,
			reflect.Uint8,
			reflect.Uint:
			v, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				fmt.Printf("parser value for %s wrong, error -> %v\n", envKey, err)
				vField.SetUint(0)
				break
			}
			vField.SetUint(v)
			break
		case reflect.String:
			vField.SetString(envValue)
			break
		default:
			break
		}
	}
	return nil
}
