package sub

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func StringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}

	return false
}

func ReadCSVFromFile(name string, separator string, fieldCount int) ([][]string, error) {
	ok, err := FileExists(name)
	if err != nil || !ok {
		return [][]string{}, err
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return [][]string{}, err
	}

	rows := strings.Split(string(data), "\n")
	var result [][]string

	for _, row := range rows {
		if strings.Count(row, separator) != fieldCount-1 {
			continue
		}
		result = append(result, strings.Split(row, separator))
	}

	return result, nil
}

func ReadFromFile(name string) (string, error) {
	ok, err := FileExists(name)
	if err != nil || !ok {
		return "", err
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func CheckInts(indexes []int, sourceArray []string) ([]int, error) {
	var intData []int
	for _, i := range indexes {
		iData, err := strconv.Atoi(sourceArray[i])
		if err != nil {
			return []int{}, err
		}

		intData = append(intData, iData)
	}

	return intData, nil
}

func RecursivePrint(object interface{}, level int) {
	isInt := reflect.TypeOf(object) == reflect.TypeOf(1)
	isFloat := reflect.TypeOf(object) == reflect.TypeOf(1.0)
	isString := reflect.TypeOf(object) == reflect.TypeOf("")
	isBool := reflect.TypeOf(object) == reflect.TypeOf(true)

	if isInt || isFloat || isString || isBool {
		for i := 0; i < level; i++ {
			fmt.Print("  ")
		}
		fmt.Println(object)
		return
	}

	if reflect.TypeOf(object).Kind() == reflect.Slice {
		objReflect := reflect.ValueOf(object)

		for i := 0; i < objReflect.Len(); i++ {
			for i := 0; i < level; i++ {
				fmt.Print("  ")
			}
			fmt.Println(objReflect.Index(i).Interface())
		}

		return
	}

	// Do it for map
	if reflect.TypeOf(object).Kind() == reflect.Map {
		fmt.Println(object)

		return
	}
	objReflect := reflect.ValueOf(object)

	for i := 0; i < objReflect.NumField(); i++ {
		for i := 0; i < level; i++ {
			fmt.Print("  ")
		}
		fmt.Println(objReflect.Type().Field(i).Name)

		RecursivePrint(objReflect.Field(i).Interface(), level+1)
	}
}
