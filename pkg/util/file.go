package util

import "io/ioutil"

// ReadFile reads the content of a file, returns it as byte slice
func ReadFile(file string) (data []byte) {
	data, _ = ioutil.ReadFile(file)
	return
}

// ReadFileASString reads the content of a file, returns it as string
func ReadFileASString(file string) string {
	return string(ReadFile(file))
}
