package record

import "fmt"

type Record struct {
    Id int 
    Cents int
    Name string 
    Category string
}

func (r *Record) String() string {
    return fmt.Sprintf("%v: id=%v cents=%v [%v]", r.Name, r.Id, r.Cents, r.Category) 
}
