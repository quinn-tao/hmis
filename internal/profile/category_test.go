package profile_test

import (
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/profile"
)

func TestCategoryAdd(t *testing.T) {
    tcs := []struct {
        Name string 
        NewPath string 
        Success bool
    } {
        {
            Name: "[TestCategoryAdd] Test add as single child", 
            NewPath: "t1/new",
            Success: true,
        }, 
        {
            Name: "[TestCategoryAdd] Test add as a sibling", 
            NewPath: "t2/t4/new",
            Success: true,
        }, 
        {
            Name: "[TestCategoryAdd] Test add failed as non-existent path", 
            NewPath: "t8/new",
            Success: false,
        }, {
            Name: "[TestCategoryAdd] Test add failed as duplicate category", 
            NewPath: "t1",
            Success: false,
        }, {
            Name: "[TestCategoryAdd] Test add at root level",
            NewPath: "t0",
            Success: true,
        },
    }

    prepComplexCategories()
    
    for _, tc := range tcs {
        t.Logf("[TestCategoryAdd] running %v", tc.Name)
        
        newCategory, err := t0.AddCategory(tc.NewPath)
        if tc.Success {
            if err != nil {
                t.Fatalf("Expected no err, got %v", err)
            }
            foundCategory, exists := t0.FindCategoryWithPath(tc.NewPath)
            if !exists || !foundCategory.Equals(newCategory) {
                t.Fatalf("Expected new category be found, but didn't")
            }
        } else {
            if err == nil {
                t.Fatalf("Expected err, got none")
            }
        }
    }
}

func TestCategoryEqual(t *testing.T) {
    categoryNoSub1 := profile.Category {
        Name: "no sub 1",
    }
    categoryNoSub1Dup := profile.Category {
        Name: "no sub 1",
    }
    categoryNoSub2 := profile.Category {
        Name: "no sub 2",
    }
    categoryNoSub2Dup := profile.Category {
        Name: "no sub 2",
    }
    
    tcs := []struct {
        Name string 
        This *profile.Category
        That *profile.Category
        ExpEqual bool
    } {
        {
            Name: "two category equal",
            This: &profile.Category{
                Name: "a",
                Sub: map[string]*profile.Category{
                    categoryNoSub1.Name: &categoryNoSub1,
                    categoryNoSub2.Name: &categoryNoSub2,
                },
            },
            That: &profile.Category{
                Name: "a",
                Sub: map[string]*profile.Category{
                    categoryNoSub1Dup.Name: &categoryNoSub1Dup,
                    categoryNoSub2Dup.Name: &categoryNoSub2Dup,
                },
            },
            ExpEqual: true,
        },
        {
            Name: "two category sub equal but name differ",
            This: &profile.Category{
                Name: "a",
                Sub: map[string]*profile.Category{
                    categoryNoSub1.Name: &categoryNoSub1,
                    categoryNoSub2.Name: &categoryNoSub2,
                },
            },
            That: &profile.Category{
                Name: "b",
                Sub: map[string]*profile.Category{
                    categoryNoSub1Dup.Name: &categoryNoSub1Dup,
                    categoryNoSub2Dup.Name: &categoryNoSub2Dup,
                },
            },
            ExpEqual: false,
        },
        {
            Name: "two category name equal but sub differ",
            This: &profile.Category{
                Name: "a",
                Sub: map[string]*profile.Category{
                    categoryNoSub1.Name: &categoryNoSub1,
                    categoryNoSub2.Name: &categoryNoSub2,
                },
            },
            That: &profile.Category{
                Name: "a",
                Sub: map[string]*profile.Category{
                    categoryNoSub2Dup.Name: &categoryNoSub2Dup,
                },
            },
            ExpEqual: false,
        },
    }

    for _, tc := range tcs {
        t.Logf("[TestCategoryEqual] Running %v", tc.Name)
        actualEqual := tc.This.Equals(tc.That)
        if actualEqual != tc.ExpEqual {
            t.Fatalf("Expected %v, Got %v", tc.ExpEqual, actualEqual)
        }
    }
}

func TestFindCategoryWithPath(t *testing.T) {
    prepComplexCategories()

    tcs := []struct {
        Name string 
        TargetPath string 
        ExpFound bool
        ExpCT *profile.Category
    } { 
        {
            Name: "one node found", 
            TargetPath: "t1",
            ExpFound: true, 
            ExpCT: &t1,
        },
        {
            Name: "one node not found - not matching", 
            TargetPath: "t3",
            ExpFound: false, 
            ExpCT: nil,
        }, {
            Name: "one node not found - path too long", 
            TargetPath: "t1/t9",
            ExpFound: false, 
            ExpCT: nil,
        },{
            Name: "multi node found", 
            TargetPath: "t2/t3",
            ExpFound: true, 
            ExpCT: &t3,
        },
        {
            Name: "multi node found - down the path", 
            TargetPath: "t2/t3/t6",
            ExpFound: true, 
            ExpCT: &t6,
        }, {
            Name: "one node not found - node exists but not on the path", 
            TargetPath: "t2/t6",
            ExpFound: false, 
            ExpCT: nil,
        },
    }
    

    for _, tc := range tcs {
        t.Logf("[TestFindCategoryWithPath] Running %v", tc.Name)
        actualCT, actualFound := t0.FindCategoryWithPath(tc.TargetPath)
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

    prepSimpleCategories()

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

func prepSimpleCategories() {
    // Form a tree:
    // multi ---> single 
    //      |
    //       ---> node1 ---> node2
    MultiNode.Sub["single"] = SingleNode  
    MultiNode.Sub["node1"] = ChainNode1
    ChainNode1.Sub["node2"] = ChainNode2
}

var (
    t0 profile.Category
    t1 profile.Category
    t2 profile.Category
    t3 profile.Category
    t4 profile.Category
    t5 profile.Category
    t6 profile.Category
    t7 profile.Category
    t8 profile.Category
    t9 profile.Category
)

func prepComplexCategories() {
    // Form a tree:
    //  t0 --- t1 
    //     | 
    //      --- t2 --- t3 --- t6
    //            | 
    //             --- t4 --- t7 --- t8
    //            |             | 
    //             --- t5        --- t9 

    t0.Name = "t0"
    t1.Name = "t1"
    t2.Name = "t2"
    t3.Name = "t3"
    t4.Name = "t4"
    t5.Name = "t5"
    t6.Name = "t6"
    t7.Name = "t7"
    t8.Name = "t8"
    t9.Name = "t9" 

    t0.Sub = map[string]*profile.Category{t1.Name:&t1, t2.Name:&t2,}
    t2.Sub = map[string]*profile.Category{
        t3.Name: &t3, t4.Name: &t4, t5.Name: &t5,
    }
    t3.Sub = map[string]*profile.Category{t6.Name: &t6}
    t4.Sub = map[string]*profile.Category{t6.Name: &t7}
    t7.Sub = map[string]*profile.Category{t6.Name: &t8, t9.Name: &t9}
}


