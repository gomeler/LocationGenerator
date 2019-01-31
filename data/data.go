package data


import (
	"os"
	"bufio"
	"encoding/gob"
)


func Load(filename string) ([]string, error) {
	data, ingestError := loadFile(filename)
	return data, ingestError
}

func loadFile(filename string) ([]string, error) {
	//given a filename, open that file, load it, close the file.
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	data, ingestError := ingestFile(file)
	return data, ingestError
}

func ingestFile(file *os.File) ([]string, error) {
	//Dump the given file line by line into a string array.
	scanner := bufio.NewScanner(file)
	dataDump := make([]string, 0)
	for scanner.Scan() {
		dataDump = append(dataDump, scanner.Text())
	}
	err := scanner.Err()
	return dataDump, err
}

func WriteGob(filename string, object interface{}) error {
	file, err := os.Create(filename)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	return err
}

func ReadGob(filename string, object interface{}) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	return err
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
