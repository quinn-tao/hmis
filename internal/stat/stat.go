package stat

import (
	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/profile"
)

type Stat struct {
    fixedExp coins.RawAmountVal
    recordedExp coins.RawAmountVal
    profile *profile.Profile
}

func NewStat() (* Stat, error) {
    var st Stat 
    currProfile, err := profile.GetCurrentProfile()
    if err != nil {
        return nil, err
    }
    st.profile = currProfile
    
    st.fixedExp = st.calcFixedExpenses()
    recordedExp, err := st.calcRecordExpenses()
    if err != nil {
        return nil, err
    }
    st.recordedExp = recordedExp

    return &st, nil
}

func (st *Stat) calcFixedExpenses() coins.RawAmountVal {
    root := st.profile.Category
    _, acc := root.Visit(profile.CategorySelectAll, func(c *profile.Category, i interface{}) interface{} {
        if c.Recurr == nil {
            return i
        }
        acc := i.(int64)
        acc += int64(c.Recurr.Amount)
        return acc
    }, int64(0))
    
    return coins.RawAmountVal(acc.(int64))
}

func (st *Stat) calcRecordExpenses() (coins.RawAmountVal, error) {
    sumRec, err := db.GetSumRecord()
    if err != nil {
        return coins.InvalidRawAmountVal, err
    }
    
    return sumRec.Amount, nil
}

func (st *Stat) GetFixedExp() coins.RawAmountVal {
    return st.fixedExp
}

func (st *Stat) GetRecordedExp() coins.RawAmountVal {
    return st.recordedExp
}

func (st *Stat) GetTotalExp() coins.RawAmountVal {
    return st.recordedExp + st.fixedExp
}

func (st *Stat) GetRemainingBudget() coins.RawAmountVal {
    return st.profile.Limit - st.recordedExp - st.fixedExp
}

func (st *Stat) GetMode() profile.Mode {
    return st.profile.Mode
}

func (st *Stat) GetBudget() coins.RawAmountVal {
    return st.profile.Limit
}
