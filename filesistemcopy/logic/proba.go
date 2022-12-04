package logic

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Map_value1 struct {
	Values     []int
	Free_space []int
}
type Item1 struct {
	Name    string
	Content *string
	Index   int
	Tag     []string
}
type GoMap struct {
	Baza          map[string]Map_value1
	Item_list     []Item1
	MainFreeSpace []int
}

func (m *GoMap) ReadTag(tag string) []Item1 {
	b := m.Baza[tag].Values
	a := []Item1{}
	//c:=Baza[tag].Free_space
	for i := 0; i < len(b); i++ {
		if m.Item_list[b[i]].Content != nil {
			a = append(a, m.Item_list[b[i]])
		}
	}
	return a
}

func (m *GoMap) Delete(index int) {
	m.MainFreeSpace = append(m.MainFreeSpace, index)
	keys := m.Item_list[index].Tag //lista stringova
	//trebamo izbrisati item
	//prvo prolazimo kroz sve tagove
	//prolzimo kroz values od itema koji brisemo
	//u free space ubacujemo slobodan index iz value
	for num := range keys {
		//go through all tags of item at given inde
		a_key := keys[num]
		values := m.Baza[a_key].Values
		for y := range values {
			if values[y] == index {
				//u values postoji neiskoristen index
				//u free space stavljamo y
				x := m.Baza[a_key]
				x.Free_space = append(x.Free_space, y)
				m.Baza[keys[num]] = x
			}
		}
	}
	stringy := m.ListToString()
	m.Update(stringy, "templates/data.txt")
	m.Item_list[index].Content = nil //free up space of deleted item
	m.Item_list[index].Name = ""
	m.Item_list[index].Tag = nil
}
func (m *GoMap) Update(WriteString string, path string) {
	if err := os.Truncate(path, 0); err != nil {
		log.Printf("Failed to truncyy &v", err)
	}
	toFile(WriteString, path)

}
func toFile(stringToWrite string, path string) {

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(stringToWrite); err != nil {
		panic(err)
	}

}
func (m *GoMap) fileToList(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m.Item_list = nil
	m.Baza = make(map[string]Map_value1)
	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	a := Item1{}
	for {

		a = Item1{}
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		a.Name = rec[0]
		a.Content = &rec[1]
		a.Index, err = strconv.Atoi(rec[2])

		for i := 3; i < len(rec); i++ {
			a.Tag = append(a.Tag, rec[i])

		}

		//logic.Item_list = append(logic.Item_list, a)
		m.Add(a)
	}

}

func (m *GoMap) ListToString() string {
	x := ""
	i := 0
	for num := range m.Item_list {
		if len(m.MainFreeSpace) == 0 {
			continue
		}
		if num == m.MainFreeSpace[i] {
			i++
			continue
		}
		x += m.Item_list[num].Name
		x += ","
		x += *(m.Item_list[num].Content)
		x += ","
		x += strconv.Itoa(m.Item_list[num].Index)
		x += ","
		for num1, val1 := range m.Item_list[num].Tag {
			x += m.Item_list[num].Tag[num1]
			if num1 != len(val1) {
				x += ","
			}
		}
		x += "\n"

	}
	return x
}
func (m *GoMap) listItemToString(parsedItem Item1) string {
	x := ""
	x += parsedItem.Name
	x += ","
	x += *parsedItem.Content
	x += ","
	x += strconv.Itoa(parsedItem.Index)
	x += ","
	for num, val := range parsedItem.Tag {
		x += parsedItem.Tag[num]
		if num != len(val)-1 {
			x += ","
		}
	}
	return x
}

func (m *GoMap) CreateTag(tagName string) {
	//ukoliko tag vec postoji
	_, ok := m.Baza[tagName]
	if ok == true {
		return
	}
	//ukoliko tag ne postoji
	fmt.Println(tagName)
	if m.Baza == nil {
		m.Baza = make(map[string]Map_value1)
	}
	m.Baza[tagName] = Map_value1{}
}

/*createItem*/
func (m *GoMap) Add(newItem Item1) {
	// kao index newItema se uvijek salje nula
	//jer ne mozemo znati na kojem se mjestu nalazi
	//dodajemo item u glavni array
	//u prvo slobodno mjesto ako postoji
	if len(m.MainFreeSpace) != 0 {
		indexIzFreeSpace := m.MainFreeSpace[len(m.MainFreeSpace)-1]
		//mainFreeSpace popback
		if len(m.MainFreeSpace) > 0 {
			m.MainFreeSpace = m.MainFreeSpace[:len(m.MainFreeSpace)-1]
		}
		m.Item_list[indexIzFreeSpace] = newItem
		newItem.Index = indexIzFreeSpace
	} else {
		//normalno stavljamo item u item list
		newItem.Index = len(m.Item_list)
		m.Item_list = append(m.Item_list, newItem)
	}

	//dodajemo item tamo gdje pripada po tagu
	for i := range newItem.Tag {
		a_tag := newItem.Tag[i]
		//ukoliko tag ne postoji kreiramo tag
		_, ok := m.Baza[a_tag]
		if ok != true {
			m.CreateTag(a_tag)
		}

		//ako tag postoji
		var tag_koji_postoji = m.Baza[a_tag]
		//ukoliko ne postoji ni jedno slobodno mjesto dodajemo na kraj
		if len(tag_koji_postoji.Free_space) == 0 {
			x := m.Baza[a_tag]
			x.Values = append(m.Baza[a_tag].Values, newItem.Index)
			m.Baza[a_tag] = x
		} else {
			//ukoliko postoji slobodno mjesto
			y := m.Baza[a_tag]
			y.Values[y.Free_space[len(y.Free_space)-1]] = newItem.Index

			if len(y.Free_space) > 0 {
				y.Free_space = y.Free_space[:len(y.Free_space)-1]
			}
			m.Baza[a_tag] = y
		}
	}
	stringy := m.listItemToString(newItem) + "\n"
	toFile(stringy, "templates/data.txt")

}
