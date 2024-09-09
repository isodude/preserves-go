# preserves-go
A Go implementation of https://git.syndicate-lang.org/syndicate-lang/preserves


# Steps to integration

## Parsing preserves, binary
## Printing preserves, binary
## Parsing preserves schema binary
## Printing preserves schema binary
### Find matches between preserves structures and host language structures
#### Union
A union in preserves is parsed as 
```
Structure = First / Second
First = int
Second = bool
```
This is then parsed as preserves
```
<or First Second>
```

In the compiled parser it should first check if the value matches First and then Second. Which means
```
5 == Structure[First]
true == Structure[Second]
```

In Go there's two ways to describe it:
```
type Structure interface {
    IsStructure()
}
type First int
type Second bool
func (_ First) IsStructure() {}
type (_ Second) IsStructure() {}
```

```
type StructureType int
const (
    StructureNone StructureType = iota
    StructureFirst
    StructureSecond
)
type Structure struct {
    Typ StructureType
    Value any
}
```

It seems that the first one is more clean.

Bolting on more needed functions
```
func NewFirst(i int) Structure {
    j := First(i)
    return &j
}
func NewSecond(b bool) Structure {
    a := Second(b)
    return &a
}

func ParseStructure(data []byte) (Structure, error) {
    rFirst := regexp.MustCompile("^[0-9]+$")
    rSecond := regexp.MustCompile("^#(t|f)$")
    if rFirst.Match(data) {
        i, err := strconv.Atoi(string(data))
        if err != nil {
            return nil, err
        }
        f := First(i)
        return &f
    }
    if rSecond.Match(data) {
        s := string(data)
        var b bool
        if s[1] == "t" {
            b = true
        }
        s := Second(b)
        return &s
    }
    return nil, fmt.Errorf("Unable to parse Structure")
}
```

#### Embedded
#### Sequence (seqof, [])
```
<seqof @pattern SimplePattern>
```

```
type Sequence struct {
    Pattern []Type
}
```
#### Set
```
type Set struct {
    Pattern map[Type]struct{}
}
```
#### Dict
```
type Dict map[Key]Value
```
#### Rec
```
type Name struct {
    Label <type>
    Fields <type>
}
```
#### Tuple
```
type Tuple struct {
    Patterns []Type
}
```
#### TuplePrefix
```
type TuplePrefix struct {
    Fixed []Type
    Variable Type
}
```
#### Atoms
Atoms are the basic types, strings, ints etc. In go it's pretty easy
```
import (
    "math"
)
type Boolean bool
type Double float64
type SignedInteger math.Big
type String string
type ByteString string
type Symbol string
```

`string` is very easy converted to `[]byte` which is not Comparable and can't be used as keys in maps, hence avoid it as an atom.