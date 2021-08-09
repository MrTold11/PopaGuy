package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

func NewFileStorageAnek() (fileStorageAnek, error) {
	FileStorageAnek	:= fileStorageAnek{

	}
	return FileStorageAnek, nil
}

func (s *InmemoryStorageAnek) getRandomAnek() (Anek, error) 	{
	fmt.Println("SECOND PART")
	c := Anek{
		name : "FUCK",
		text : "THIS",
	}
	str, e := os.ReadFile("number2.txt")
	if e != nil {
		fmt.Println(e)
		return c, e
	}
	st := string (str)
	number, er := strconv.Atoi(st)
	if er != nil {
		fmt.Println(er)
		return c, er
	}
	index := rand.Intn(number)
	if index == 0 {
		index++
	}
	fileName := fmt.Sprintf("anekdot%d", index)
	data, err := ioutil.ReadFile(fileName + ".json")
	fmt.Println(fileName)
	fmt.Println(string (data) )
	var a Anek
		err = json.Unmarshal(data, &a)
		if err != nil {
			fmt.Println("TRY")
			fmt.Println(err)
			b := c
			return b, err
		} else {
			fmt.Println("PLEASE WORK")
			fmt.Println(a.text)
			return a, nil
		}


}

// get random  - redo to work with files
func (s *fileStorageAnek) getRandomAnek() (Anek, error) {
	fmt.Println("WE ARE IN FIRST PART")
	data, err := s.storage.getRandomAnek()
	if err != nil {
		fmt.Println(err)
		return data,err
	} else {
		fmt.Println(data.text)
		return data, err
	}
}