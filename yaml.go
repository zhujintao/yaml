package yaml

import (
	"io/ioutil"
	"strings"
	"fmt"


)

type yaml struct {
	lines 	[]string
	cpath 	string
	num 	int

	ktbl 	map[string]int
}


func NewReadYaml(filepath string) *yaml {
	bFile,e:=ioutil.ReadFile(filepath)
	if e!=nil{  return nil }


	return &yaml{
		lines:strings.Split(string(bFile),"\n"),
		cpath:filepath,
		num:-1,

	}

}


func (y *yaml) find(key string) error {

	keys:=strings.Split(key,":")

	keyindex:=0

	y.ktbl=make(map[string]int)
	for i,line:=range y.lines {


		if len(strings.Split(line,"#")[0]) == 0 {continue}
		if len(strings.Fields(line)) == 0 {continue}

		if len(keys) == 2 {


			if y.ktbl[keys[0]] != 0 {

				if !strings.Contains(line," "){
					y.ktbl["#done"]=i
					return nil
				}
			}

			if strings.Contains(line, keys[0] + ":"){
				y.ktbl[keys[0]]=i
				y.num=i
			}

			continue

		}


		if keyindex < len(keys)-1 {

			if strings.Contains(line, keys[keyindex] + ":") {

				l:=strings.Split(line,":")[0]
				if keys[keyindex] != strings.Fields(l)[0] {continue}

				y.ktbl[keys[keyindex]]=i
				keyindex++

				if keyindex == len(keys)-1{
					y.num = i
					return nil
				}
			}
		}

	}
	return fmt.Errorf(`KEY: "%s" not found.` + "\n",key)

}


func (y *yaml) Set(key, value string) error {

	if e:=y.find(key);e!=nil {
		//fmt.Printf(e.Error())
		return e

	}

	old:=strings.Split(y.lines[y.num],":")
	k:=string([]rune(y.lines[y.num])[0:len(old[0])+1])

	y.lines[y.num]= k + " " + value

	return nil

}
func (y *yaml) Save(t ...bool) error {
	if t == nil {

		y.print()
		return nil
	}

	if t[0] {

		return ioutil.WriteFile(y.cpath,[]byte(strings.Join(y.lines,"\n")),0644)
	}
	return nil
}


func (y *yaml) Add(key, value string) error {

	if e:=y.find(key);e==nil{
		//fmt.Printf(`KEY: "%s" is exist.` + "\n",key)
		return e
	}
	keys:=strings.Split(key,":")





	if y.ktbl[keys[0]] == 0 {
		whitespace:="  "
		fkey:=keys[0] + ":\n"
		skey:=""

		for i,key:=range keys[1:] {

			if i == len(keys[1:]) -2 {
				skey+=whitespace + key + ":"
				skey+=" " + value
				break
			}

			skey+=whitespace + key + ":\n"
			whitespace+=whitespace
		}



		tmpline:=fmt.Sprintf(fkey + skey +"\n")
		y.lines = append(y.lines,tmpline)
	}else {

		whitespace:="  "
		skey:=""

		for range keys[1:len(keys)-2] {

			whitespace+=whitespace
		}

		skey+=whitespace + keys[len(keys)-2] + ":" + " " + value

		tmplines:=[]string{}

		tmpline:=fmt.Sprintf(skey +"\n")

		tmplines=append(tmplines,y.lines[0:y.ktbl[keys[len(keys)-3]]+1]...)
		tmplines=append(tmplines,tmpline)
		tmplines=append(tmplines,y.lines[y.ktbl[keys[len(keys)-3]]+1:]...)
		y.lines=tmplines
	}
	return nil

}


func (y *yaml) Get(key string) string {

	if e:=y.find(key);e!=nil{
		fmt.Printf(e.Error())
		return ""
	}
	keys:=strings.Split(key,":")
	if len(keys)-1 == 1 {return "" }

	old:=strings.Split(y.lines[y.num],keys[len(keys)-1]+":")


	v:=strings.Split(old[1],"#")[0]

	return strings.Fields(v)[0]

}

func (y *yaml) Del(key string) error {

	if e:=y.find(key);e!=nil{
		//fmt.Printf(`KEY: "%s" not found.` + "\n",key)
		return e
	}

	keys:=strings.Split(key,":")

	if len(keys) == 2 {

		y.lines=append(y.lines[0:y.ktbl[keys[0]]], y.lines[y.ktbl["#done"]-1:]...)
	}


	y.lines=append(y.lines[0:y.num], y.lines[y.num+1:]...)

	//y.lines[y.num]=""

	return nil

}

func (y *yaml) SetA(key,value string) error  {

	if e:=y.Set(key, value); e!=nil {

		return y.Add(key, value)

	}

	return nil

}


func (y *yaml) print() {
	fmt.Printf(strings.Join(y.lines,"\n"))

}

