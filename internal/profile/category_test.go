package profile_test

import (
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/profile"
)

func TestFindCategoryWithPath(t *testing.T) {
    tcs := []struct {
        Name string 
        TargetPath string 
        CT *profile.Category 
        ExpFound bool
        ExpCT *profile.Category
    } { 
        {
            Name: "one node found", 
            TargetPath: "single",
            CT: SingleNode,
            ExpFound: true, 
            ExpCT: SingleNode,
        },
        {
            Name: "one node not found - not matching", 
            TargetPath: "default",
            CT: SingleNode,
            ExpFound: false, 
            ExpCT: nil,
        }, {
            Name: "one node not found - path too long", 
            TargetPath: "single/default",
            CT: SingleNode,
            ExpFound: false, 
            ExpCT: nil,
        },{
            Name: "multi node found", 
            TargetPath: "multi/node1",
            CT: MultiNode,
            ExpFound: true, 
            ExpCT: ChainNode1,
        },
        {
            Name: "multi node found - down the path", 
            TargetPath: "multi/node1/node2",
            CT: MultiNode,
            ExpFound: true, 
            ExpCT: ChainNode2,
        }, {
            Name: "one node not found - node exists but not on the path", 
            TargetPath: "multi/node2",
            CT: MultiNode,
            ExpFound: false, 
            ExpCT: nil,
        },
    }

    prepareCategoryTrees()

    for _, tc := range tcs {
        t.Logf("[TestFindCategoryWithPath] Running %v", tc.Name)
        actualCT, actualFound := tc.CT.FindCategoryWithPath(tc.TargetPath)
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

var ChainNode1 = &profile.Category {
    Name: "node1",
    Sub: make(map[string]*profile.Category, 0),
}

var ChainNode2 = &profile.Category {
    Name: "node2",
}

func prepareCategoryTrees() {
    // Form a tree:
    // multi ---> single 
    //      |
    //       ---> node1 ---> node2
    MultiNode.Sub["single"] = SingleNode  
    MultiNode.Sub["node1"] = ChainNode1
    ChainNode1.Sub["node2"] = ChainNode2
}
