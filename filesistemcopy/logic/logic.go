package logic

/*
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Baza = make(map[string]Map_value)
var Item_list []Item
var MainFreeSpace []int

type Map_value struct {
	Values     []int
	Free_space []int
}
type Item struct {
	Name    string
	Content string
	Index   int
	Tag     []string
}

func GetItemByIndex(index int) *Item {
	return &Item_list[index]
}

func Read() {
	for _, v := range Item_list {
		if v.Content != "nil" {
			fmt.Print("Name: ", v.Name, "\n")
			fmt.Print("Content: ", v.Content, "\n")
			fmt.Print("Tags: ", v.Tag, "\n\n")
		}
	}

}
func ReadTag(tag string) []Item {

	b := Baza[tag].Values
	a := []Item{}
	//c:=Baza[tag].Free_space
	for i := 0; i < len(b); i++ {
		if Item_list[b[i]].Content != "nil" {
			a = append(a, Item_list[b[i]])
			fmt.Print("Name: ", Item_list[b[i]].Name, "\n")
			fmt.Print("Content: ", Item_list[b[i]].Content, "\n")
			fmt.Print("Tags: ", Item_list[b[i]].Tag, "\n")
		}
	}
	return a
}

func DeleteItem(index int) {

	MainFreeSpace = append(MainFreeSpace, index)
	keys := Item_list[index].Tag //lista stringova
	//trebamo izbrisati item
	//prvo prolazimo kroz sve tagove
	//prolzimo kroz values od itema koji brisemo
	//u free space ubacujemo slobodan index iz value
	for num := range keys {
		//go through all tags of item at given inde
		a_key := keys[num]
		values := Baza[a_key].Values
		for y := range values {
			if values[y] == index {
				//u values postoji neiskoristen index
				//u free space stavljamo y
				x := Baza[a_key]
				x.Free_space = append(x.Free_space, y)
				Baza[keys[num]] = x
			}
		}
	}

	Item_list[index].Content = "nil" //free up space of deleted item
	Item_list[index].Name = ""
	Item_list[index].Tag = nil

}

func (i *Item) Edit() {
	var choice string
	fmt.Println("Izaberite koji od podatak zelite promijeniti: \n 1) Ime teme \n 2) Sadrzaj \n 3)Tag")
	fmt.Scan(&choice)

	for choice != "Ime" && choice != "Sadrzaj" && choice != "1" && choice != "2" && choice != "ime" && choice != "sadrzaj" && choice != "Tag" && choice != "tag" && choice != "3" {
		fmt.Println("Neispravan unos pokušajte ponovo \n NAPOMENA: unos moze biti u obliku broja ili naziva ponudjene opcije")
		fmt.Scan(&choice)
	}

	if choice == "1" || choice == "ime" {
		choice = "Ime"
	}
	if choice == "2" || choice == "sadrzaj" {
		choice = "Sadrzaj"
	}
	if choice == "3" || choice == "tag" {
		choice = "Tag"

	}

	switch choice {
	case "Ime":
		i.Edit_name()
	case "Sadrzaj":
		i.Edit_content()
	case "Tag":
		i.Edit_tag()
	default:
		return
	}
	return

}
func (x *Item) Edit_name() {
	var new_name string
	fmt.Println("Kako zelite da se nova tema zove? ")
	fmt.Scan(&new_name)
	x.Name = new_name
	Item_list[x.Index] = *x
	return

}

func (x *Item) Edit_content() {
	var new_content string
	fmt.Println("Stari sadrzaj:\n ", x.Content)
	fmt.Println("Unesite novi sadrzaj teme: ")

	inputReader := bufio.NewReader(os.Stdin)
	new_content, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	new_content = strings.TrimSuffix(new_content, "\n")
	x.Content = new_content
	Item_list[x.Index] = *x
	return
}

func (x *Item) Edit_tag() {
	var new_tag string
	fmt.Println("Koji tag zelite dodati temi")
	fmt.Scan(&new_tag)
	x.Tag = append(x.Tag, new_tag)
	_, ok := Baza[new_tag]
	if ok == true {
		x := Map_value{Baza[new_tag].Free_space, append(Baza[new_tag].Values, x.Index)}
		Baza[new_tag] = x
	} else {
		Baza[new_tag] = Map_value{[]int{x.Index}, []int{}}
	}
	Item_list[x.Index] = *x
	return
}

/* var baza = make(map[string]map_value)
 var item_list []item
var i *template.Template = template.Must(template.ParseFiles("./templates/index.html"))
var s *template.Template = template.Must(template.ParseFiles("./templates/select.html"))
var a *template.Template = template.Must(template.ParseFiles("./templates/adding.html"))*/
/*
func CreateTag(tagName string) {
	//ukoliko tag vec postoji
	_, ok := Baza[tagName]
	if ok == true {
		fmt.Println("Tag already exists")
	}
	//ukoliko tag ne postoji
	Baza[tagName] = Map_value{}
}

func CreateItem(newItem Item) {
	// kao index newItema se uvijek salje nula
	//jer ne mozemo znati na kojem se mjestu nalazi
	//dodajemo item u glavni array
	//u prvo slobodno mjesto ako postoji
	if len(MainFreeSpace) != 0 {
		indexIzFreeSpace := MainFreeSpace[len(MainFreeSpace)-1]
		//mainFreeSpace popback
		if len(MainFreeSpace) > 0 {
			MainFreeSpace = MainFreeSpace[:len(MainFreeSpace)-1]
		}
		fmt.Println(MainFreeSpace)
		Item_list[indexIzFreeSpace] = newItem
		newItem.Index = indexIzFreeSpace
	} else {
		//normalno stavljamo item u item list
		newItem.Index = len(Item_list)
		Item_list = append(Item_list, newItem)
	}

	//dodajemo item tamo gdje pripada po tagu
	for i := range newItem.Tag {
		a_tag := newItem.Tag[i]
		//ukoliko tag ne postoji kreiramo tag
		_, ok := Baza[a_tag]
		if ok != true {
			CreateTag(a_tag)
		}

		//ako tag postoji
		var tag_koji_postoji = Baza[a_tag]
		//ukoliko ne postoji ni jedno slobodno mjesto dodajemo na kraj
		if len(tag_koji_postoji.Free_space) == 0 {
			x := Baza[a_tag]
			x.Values = append(Baza[a_tag].Values, newItem.Index)
			Baza[a_tag] = x
		} else {
			//ukoliko postoji slobodno mjesto
			y := Baza[a_tag]
			y.Values[y.Free_space[len(y.Free_space)-1]] = newItem.Index

			if len(y.Free_space) > 0 {
				y.Free_space = y.Free_space[:len(y.Free_space)-1]
			}
			Baza[a_tag] = y
		}
	}
}
*/
