package profile

import (
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/quinn-tao/hmis/v1/config"
)

type Profile struct {
    Name     string
    Profiles []string

    // Loadable from profile 
    Budget uint32
}

// The current loaded profile 
var currProfile *Profile

func LoadProfile() error {
    var newProfile Profile
    dir := config.ProfileDir()

    profiles, err := os.ReadDir(dir) 
    if err != nil {
        log.Panicf("Error loading profile: %v", err)
        return err
    }
    newProfile.Profiles = make([]string, 0)
    
    name := config.CurrentProfileName()
    for _, profile := range profiles {
        if profile.Name() == name {
            newProfile.Name = profile.Name()
        }
        newProfile.Profiles = append(newProfile.Profiles, profile.Name())
    } 

    currProfile = &newProfile
    return nil
}

func Dump() {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.Style().Options.SeparateRows = false
    t.Style().Options.SeparateColumns = false
    t.Style().Options.DrawBorder = false

    t.AppendRows([]table.Row{
        {"Current Profile:", currProfile.Name,},
    })
    t.AppendSeparator()
    t.AppendRows([]table.Row{
        {"List of Available Profiles:"},
    })
    
    for _, profile := range currProfile.Profiles {
        t.AppendRows([]table.Row{
            {profile},
        })
    }

    t.Render()
}


