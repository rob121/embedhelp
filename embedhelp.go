package embedhelp

import (
	"embed"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

/*
This package supports taking embedded items in a go binary and dumping them to the os.
Perhaps you want to pack up a service file or dump the html and serve locally, this will facilitate that

Usage:
//go:embed example.service
var svcvar string
//go:embed example
var sys embed.FS
//go:embed test.txt
var test string
embedhelp.Register("service",svcvar,"/var/lib/systemd/example.service",false)
embedhelp.Register("example",sys,"/tmp/embedtest",true)
embedhelp.Register("example2",test,"/tmp/test.txt",true)
embedhelp.DumpAll()
*/


var items map[string]item

type item struct{
	contents interface{}
	dest string
	force bool
}


func init(){
	items = make(map[string]item)
}


func Register(key string,contents interface{},dest string,force bool){

	items[key]=item{contents,dest,force}

}

func DumpItem(key string) error{

    if it,ok := items[key]; ok {

    	return dump(it)

	}

	return errors.New("Item not found")

}

func DumpAll() []error{
    var errs []error
	for _,it := range items {

	  err := dump(it)

	  if(err!=nil){
	  	errs = append(errs,err)
	  }

	}
    return errs
}

func dump(it item) error{

    typ := reflect.TypeOf(it.contents).String()

	switch typ {
	case "string":
		 return dumpstr(it)
	case "embed.FS":
	     return dumpfs(it)
	default:
		 return dumpbyte(it)
	}


    return nil


}

func dumpbyte(it item) error{

	fl := it.dest

	contents := it.contents.([]byte)

	if _, err := os.Stat(fl); err == nil {

		if it.force == true {
			//write it out

			return ioutil.WriteFile(fl, contents, 0644)

		}


	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist

		return ioutil.WriteFile(fl, contents, 0644)

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		return err

	}

	return nil


}

func dumpfs(it item) error{

	contents := it.contents.(embed.FS)

	if _, err := os.Stat(it.dest); os.IsNotExist(err) {
		merr := os.Mkdir(it.dest, 0644)
		if(it.force!=true && merr!=nil){
			return merr
		}
	}


	return fs.WalkDir(contents, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if(d.IsDir()){

			if _, err := os.Stat(filepath.Join(it.dest,path)); os.IsNotExist(err) {
				merr := os.Mkdir(filepath.Join(it.dest,path), 0644)

				if(it.force!=true && merr!=nil){

					return merr
				}
			}


		}else{

			data, rerr := fs.ReadFile(contents, path)
			if rerr != nil {
				return rerr
			}

			nitem := item{data,filepath.Join(it.dest,path),it.force}
			return dumpbyte(nitem)
		}

		return nil
	})

   return nil

}

func dumpstr(it item) error{

    fl := it.dest

   contents := it.contents.(string)

	if _, err := os.Stat(fl); err == nil {

		if it.force == true {
			//write it out

			return ioutil.WriteFile(fl, []byte(contents), 0644)

		}


	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist

		return ioutil.WriteFile(fl, []byte(contents), 0644)

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		return err

	}

    return nil

}






