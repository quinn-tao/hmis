package chrono_test

import (
	"fmt"
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/chrono"
)

func TestFormatDate(t *testing.T) {
    t.Logf("[TestFormatDate] sanity")
    date, err := chrono.NewDate("2024-01-02")
    if err != nil {
        t.Fatal("Not expecting error whilk")
    }
    
    dateFmt := fmt.Sprintf("%v", date) 
    if dateFmt != "2024-01-02" {
        t.Fatalf("Expecting 2024-01-02; Got %v", dateFmt)
    }
}
