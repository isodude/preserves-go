package beep

import (
	. "github.com/isodude/preserves-go/lib/preserves"
)

type Date struct {
	Year  SignedInteger
	Month SignedInteger
	Day   SignedInteger
}

func NewDate(year SignedInteger, month SignedInteger, day SignedInteger) *Date {
	return &Date{Year: year, Month: month, Day: day}
}
func DateFromPreserves(value Value) *Date {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 3 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("date")) {
				if p0 := SignedIntegerFromPreserves(rec.Fields[0]); p0 != nil {
					if p1 := SignedIntegerFromPreserves(rec.Fields[1]); p1 != nil {
						if p2 := SignedIntegerFromPreserves(rec.Fields[2]); p2 != nil {
							return &Date{Year: *p0, Month: *p1, Day: *p2}
						}
					}
				}
			}
		}
	}
	return nil
}
func DateToPreserves(s Date) Value {
	return &Record{Key: NewSymbol("date"), Fields: []Value{SignedIntegerToPreserves(s.Year), SignedIntegerToPreserves(s.Month), SignedIntegerToPreserves(s.Day)}}
}

type Hat struct {
	Color Pstring
}

func NewHat(color Pstring) *Hat {
	return &Hat{Color: color}
}
func HatFromPreserves(value Value) *Hat {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("hat")) {
				if p0 := PstringFromPreserves(rec.Fields[0]); p0 != nil {
					return &Hat{Color: *p0}
				}
			}
		}
	}
	return nil
}
func HatToPreserves(s Hat) Value {
	return &Record{Key: NewSymbol("hat"), Fields: []Value{PstringToPreserves(s.Color)}}
}

type Item interface {
	IsItem()
}

func ItemFromPreserves(value Value) Item {
	if o := ItemHatFromPreserves(value); o != nil {
		return o
	}
	if o := ItemJacketFromPreserves(value); o != nil {
		return o
	}
	if o := ItemShoeFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func ItemToPreserves(s Item) Value {
	switch u := s.(type) {
	case *ItemHat:
		return ItemHatToPreserves(*u)
	case *ItemJacket:
		return ItemJacketToPreserves(*u)
	case *ItemShoe:
		return ItemShoeToPreserves(*u)
	}
	return nil
}

type ItemHat struct {
	Hat
}

func NewItemHat(obj Hat) *ItemHat {
	return &ItemHat{Hat: obj}
}
func (*ItemHat) IsItem() {
}
func ItemHatFromPreserves(value Value) *ItemHat {
	if o := HatFromPreserves(value); o != nil {
		return &ItemHat{Hat: *o}
	}
	return nil
}
func ItemHatToPreserves(s ItemHat) Value {
	return HatToPreserves(s.Hat)
}

type ItemJacket struct {
	Jacket
}

func NewItemJacket(obj Jacket) *ItemJacket {
	return &ItemJacket{Jacket: obj}
}
func (*ItemJacket) IsItem() {
}
func ItemJacketFromPreserves(value Value) *ItemJacket {
	if o := JacketFromPreserves(value); o != nil {
		return &ItemJacket{Jacket: *o}
	}
	return nil
}
func ItemJacketToPreserves(s ItemJacket) Value {
	return JacketToPreserves(s.Jacket)
}

type ItemShoe struct {
	Shoe
}

func NewItemShoe(obj Shoe) *ItemShoe {
	return &ItemShoe{Shoe: obj}
}
func (*ItemShoe) IsItem() {
}
func ItemShoeFromPreserves(value Value) *ItemShoe {
	if o := ShoeFromPreserves(value); o != nil {
		return &ItemShoe{Shoe: *o}
	}
	return nil
}
func ItemShoeToPreserves(s ItemShoe) Value {
	return ShoeToPreserves(s.Shoe)
}

type Jacket struct {
	Kind Pstring
}

func NewJacket(kind Pstring) *Jacket {
	return &Jacket{Kind: kind}
}
func JacketFromPreserves(value Value) *Jacket {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("jacket")) {
				if p0 := PstringFromPreserves(rec.Fields[0]); p0 != nil {
					return &Jacket{Kind: *p0}
				}
			}
		}
	}
	return nil
}
func JacketToPreserves(s Jacket) Value {
	return &Record{Key: NewSymbol("jacket"), Fields: []Value{PstringToPreserves(s.Kind)}}
}

type Object interface {
	IsObject()
}

func ObjectFromPreserves(value Value) Object {
	if o := ObjectItemFromPreserves(value); o != nil {
		return o
	}
	if o := ObjectPersonFromPreserves(value); o != nil {
		return o
	}
	return nil
}
func ObjectToPreserves(s Object) Value {
	switch u := s.(type) {
	case *ObjectItem:
		return ObjectItemToPreserves(*u)
	case *ObjectPerson:
		return ObjectPersonToPreserves(*u)
	}
	return nil
}

type ObjectItem struct {
	Item
}

func NewObjectItem(obj Item) *ObjectItem {
	return &ObjectItem{Item: obj}
}
func (*ObjectItem) IsObject() {
}
func ObjectItemFromPreserves(value Value) *ObjectItem {
	if o := ItemFromPreserves(value); o != nil {
		return &ObjectItem{Item: o}
	}
	return nil
}
func ObjectItemToPreserves(s ObjectItem) Value {
	return ItemToPreserves(s.Item)
}

type ObjectPerson struct {
	Person
}

func NewObjectPerson(obj Person) *ObjectPerson {
	return &ObjectPerson{Person: obj}
}
func (*ObjectPerson) IsObject() {
}
func ObjectPersonFromPreserves(value Value) *ObjectPerson {
	if o := PersonFromPreserves(value); o != nil {
		return &ObjectPerson{Person: *o}
	}
	return nil
}
func ObjectPersonToPreserves(s ObjectPerson) Value {
	return PersonToPreserves(s.Person)
}

type Person struct {
	Name     Pstring
	Birthday Date
	Wears    []Object
}

func NewPerson(name Pstring, birthday Date, wears []Object) *Person {
	return &Person{Name: name, Birthday: birthday, Wears: wears}
}
func PersonFromPreserves(value Value) *Person {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 4 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("person")) {
				if p0 := PstringFromPreserves(rec.Fields[0]); p0 != nil {
					if p1 := DateFromPreserves(rec.Fields[1]); p1 != nil {
						if seq, ok := rec.Fields[2].(*Sequence); ok {
							var p2 []Object
							for _, item := range *seq {
								if itemParsed := ObjectFromPreserves(item); itemParsed != nil {
									p2 = append(p2, itemParsed)
								} else {
									return nil
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
func PersonToPreserves(s Person) Value {
	var p2 = &Sequence{}
	for _, k := range s.Wears {
		*p2 = append(*p2, ObjectToPreserves(k))
	}
	return &Record{Key: NewSymbol("person"), Fields: []Value{PstringToPreserves(s.Name), DateToPreserves(s.Birthday), p2}}
}

type Shoe struct {
	Kind Pstring
}

func NewShoe(kind Pstring) *Shoe {
	return &Shoe{Kind: kind}
}
func ShoeFromPreserves(value Value) *Shoe {
	if rec, ok := value.(*Record); ok && len(rec.Fields) == 1 {
		if v := SymbolFromPreserves(rec.Key); v != nil {
			if v.Equal(NewSymbol("shoe")) {
				if p0 := PstringFromPreserves(rec.Fields[0]); p0 != nil {
					return &Shoe{Kind: *p0}
				}
			}
		}
	}
	return nil
}
func ShoeToPreserves(s Shoe) Value {
	return &Record{Key: NewSymbol("shoe"), Fields: []Value{PstringToPreserves(s.Kind)}}
}
