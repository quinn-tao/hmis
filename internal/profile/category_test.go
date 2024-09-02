package profile_test

import (
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/profile"
)


func TestFindCategoryRecursive(t *testing.T) {
    tcs := []struct {
        Name string 
        TargetName string 
        CT *profile.Category
        ExpFound bool
        ExpCT *profile.Category
    }{
        {
            Name: "one node",
            TargetName: "single",
            CT: SingleNode,
            ExpFound: true,
            ExpCT: SingleNode,
        },{
            Name: "one node not found",
            TargetName: "multi",
            CT: SingleNode,
            ExpFound: false,
        },{
            Name: "tree found sub",
            TargetName: "single",
            CT: MultiNode,
            ExpFound: true,
            ExpCT: SingleNode,
        }, {
            Name: "tree found root",
            TargetName: "multi",
            CT: MultiNode,
            ExpFound: true,
            ExpCT: MultiNode,
        }, {
            Name: "tree not found",
            TargetName: "default",
            CT: MultiNode,
            ExpFound: false,
        },
    }

    prepareCategoryTrees()

    for _, tc := range tcs {
        t.Logf("[TestFindCategoryRecursive] Running %v", tc.Name)
        actualCT, actualFound := tc.CT.FindCategoryRecursive(tc.TargetName)
        if !tc.ExpFound {
            if actualFound {
                t.Fatal("Should not find target, but actually found")
            }
            continue
        }
        if !actualFound {
            t.Fatalf("Expected %v, actually not found", tc.ExpCT)
        }
        if actualCT != tc.ExpCT {
            t.Fatalf("Expected %v, Got %v", tc.ExpCT, actualCT)
        }
    }
}

var SingleNode = &profile.Category {
    Name: "single",
    Sub: make(map[string]*profile.Category, 0),
}

var MultiNode  = &profile.Category {
    Name: "multi",
    Sub: make(map[string]*profile.Category, 0),
}

func prepareCategoryTrees() {
    MultiNode.Sub["single"] = SingleNode  
}
