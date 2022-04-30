package reflection

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func ObjectTypePath(obj any) string {
	objType := reflect.TypeOf(obj).Elem()
	path := fmt.Sprintf("%s.%s", objType.PkgPath(), objType.Name())
	return path
}

func TypePath[T any]() string {
	var msg T
	return ObjectTypePath(msg)
}

func CreateInstance[T any]() T {
	var msg T
	ttyp := reflect.TypeOf(msg).Elem()
	instance := reflect.New(ttyp).Interface()
	return instance.(T)
}

func MethodPath(f interface{}) string {
	pointerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lastSlashIdx := strings.LastIndex(pointerName, "/")
	methodPath := pointerName[lastSlashIdx+1:]
	if methodPath[len(methodPath)-3:] == "-fm" {
		methodPath = methodPath[:len(methodPath)-3]
	}
	methodPath = strings.ReplaceAll(methodPath, ".", ":")
	methodPath = strings.ReplaceAll(methodPath, "(", "")
	methodPath = strings.ReplaceAll(methodPath, ")", "")
	methodPath = strings.ReplaceAll(methodPath, "*", "")
	return methodPath
}
