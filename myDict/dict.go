package myDict

import "errors"

type Dictionary map[string]string

var errNoWord error=errors.New("there is no word")
var errDuplicatedWord error=errors.New("there is alraedy exist word")

func (d Dictionary) Search(word string)(value string,err error)  {
	var exist bool
	value,exist=d[word]
	if !exist{
		err=errNoWord
	}
	return
}

func (d Dictionary) Add(word,definition string) error{
	_,err:=d.Search(word)
	if err!=nil{
		d[word]=definition
		return nil
	}
	return errDuplicatedWord
}