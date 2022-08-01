package fcore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/xid"
)

func StructFieldList(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		fmt.Println(val.Type().Field(i).Name)
	}
}

// Flags: Multiple String Values
type ArrayFlagString []string

func (i *ArrayFlagString) String() string {
	return fmt.Sprint(*i)
}
func (i *ArrayFlagString) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Flags: Multiple Int Values
type ArrayFlagInt []int

func (i *ArrayFlagInt) String() string {
	return fmt.Sprint(*i)
}
func (i *ArrayFlagInt) Set(value string) error {
	cValue, _ := strconv.Atoi(value)
	*i = append(*i, cValue)
	return nil
}

/*func FileExist(dirName string) bool {
	_, err := os.Stat(dirName)
	return !os.IsNotExist(err)
}*/

func CreateDir(path string) error {
	// Check if folder exists
	_, err := os.Stat(path)

	// Create directory if not exists
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModeDir|0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateFile(filename string) error {
	// Check if folder exists
	_, err := os.Stat(filename)

	// Create directory if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		defer file.Close()
	}

	return nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFileOrDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(filepath string) string {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return ""
	}

	return string(data)
}

func WriteFileLineByLine(filepath string, data []string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, line := range data {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteFileAll(filepath string, data string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	file.Sync()

	return nil
}

func CalculateScore(likes, downloads, views, comments int) int {

	score := 0

	// Likes
	score += likes * 10

	// Downloads
	score += downloads * 8

	// Views
	score += views * 5

	// Comments
	score += comments * 1

	return score
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func GenerateRedisKey(key, pattern, val string) string {
	return strings.ReplaceAll(key, pattern, val)
}

func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}

func GenerateRandUID() string {
	return xid.New().String()
}

func Convert2RedisMap(item interface{}) map[string]interface{} {
	var itemMap map[string]interface{}
	data, _ := json.Marshal(item)
	json.Unmarshal(data, &itemMap)

	return itemMap
}

/*
Pagination Query
1 2 3 4 5 6 7 8 9 [10] >
<< < 5 6 7 8 9 [10] 11 12 13 14 >
<< < 15 16 17 18 19 [20] 21 22 23 24 >
*/
func GetPaginationList(total int, limit int, paged int) (int, []string, int, int) {
	var pageNavi []string

	if total == 0 {
		return 1, pageNavi, 0, 0
	}

	prev := 0
	next := 0

	// Calculate total of pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	navi := 10
	start := 1
	finish := totalPages

	if totalPages > navi {
		finish = navi

		if paged >= totalPages {
			finish = totalPages
		}

		if totalPages > 5 && paged > 5 {
			start = Max(1, paged-4)
			finish = Min(totalPages, paged+5)
		}
	}

	if paged > 5 {
		prev = 1
	}

	if totalPages > finish {
		next = finish + 1
	}

	if start != finish {
		for i := start; i <= finish; i++ {
			page := strconv.Itoa(i)
			pageNavi = append(pageNavi, page)
		}
	}

	return totalPages, pageNavi, prev, next
}

func GetTruncateText(str string, limitter int) string {
	str = MakeSingleLineString(str)
	str = SanitizeString(str)
	if utf8.RuneCountInString(str) < limitter {
		return str
	}

	return string([]rune(str)[:limitter]) + "..."
}

func SearchTextLimit(str string, limitter int) string {
	var htmlEscaper = strings.NewReplacer(
		`/`, "",
		`\`, "",
		`>`, "",
		`>`, "",
		`%3C`, "",
		`%3E`, "",
		`&lt;`, "",
		`&gt;`, "",
		`(`, "",
		`)`, "",
		`%22`, "",
		`"`, "",
		`'`, "",
		`;`, "",
	)
	str = htmlEscaper.Replace(str)

	if utf8.RuneCountInString(str) < limitter {
		return str
	}

	return string([]rune(str)[:limitter])
}

func MakeSingleLineString(str string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(str, " ")
}

func SanitizeString(str string) string {
	re := regexp.MustCompile(`(?m)<[^>]*>`)
	return re.ReplaceAllString(str, " ")
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == string(str) {
			return true
		}
	}

	return false
}
