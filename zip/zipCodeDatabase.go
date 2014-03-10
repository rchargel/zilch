package zip

import (
	"io/ioutil"
)

type ZipCodeDB struct {
	dir string
}

func (z ZipCodeDB) LoadAll(ch chan ZipCodeEntry) {
	open := 0
	fch := make(chan ZipCodeEntry)

	files, err := ioutil.ReadDir(z.dir)
	if err != nil {
		panic("Could not read files")
	}
	for _, file := range files {
		r := ZipReader{z.dir + file.Name()}
		go r.Read(fch)
		open = open + 1
	}

	for entry := range fch {
		if len(entry.ZipCode) == 0 {
			//close connection
			open = open - 1
			if open == 0 {
				break
			}
		} else {
			ch <- entry
		}
	}
	close(ch)
}
