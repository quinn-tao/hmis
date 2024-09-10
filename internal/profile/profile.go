package profile

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"golang.org/x/text/currency"
	"gopkg.in/yaml.v3"
)

type Profile struct {
    Profiles []string
    Name string
    Limit currency.Amount
    Mode Mode
    Currency currency.Unit
    Category *Category
}

// The current loaded profile 
var currProfile *Profile

// Load current profile 
func LoadProfile() error {
    var newProfile Profile
    dir := config.ProfileDir()
    
    debug.TraceF("Loading profile directory %v", dir)
    profiles, err := os.ReadDir(dir) 
    if err != nil {
        log.Panicf("Error loading profile: %v", err)
        return err
    }
    newProfile.Profiles = make([]string, 0)
    
    name := config.CurrentProfileName()
    debug.TraceF("Config: current profile name %v", name)
    for _, file := range profiles {
        tokens := strings.Split(file.Name(),".")
        if len(tokens) < 2 || tokens[1] != "yaml" {
            // TODO: support profile names with multiple dots 
            debug.TraceF("Skipped file %v", file.Name())
            continue
        }

        profile := tokens[0]
        if profile == name {
            newProfile.Name = profile
        }
        newProfile.Profiles = append(newProfile.Profiles, profile)
    } 

    profilePath := path.Join(dir, newProfile.Name+".yaml")
    newProfile.parseProfile(profilePath)

    // Read & Parser Yaml
    debug.Trace("Profile config initialized")
    currProfile = &newProfile
    return nil
}

// Find a category 
func FindCategory(name string) (c *Category, exists bool){
    if c, exists := currProfile.Category.FindCategoryWithPath(name); exists {
        return c, exists
    }
    return currProfile.Category.FindCategoryRecursive(name)
}

// List of field parsers. Each individually load a section of 
// users's budgeting profile from generic map[interface{}]interface{}
// NOTE: order of these parsers matters as some may depends on others. 
//       For instance, recurrent expenses only make sense after currency settings 
//       are loaded
// TODO: the parsers still share a farely amount of biolerplate code 
//       maybe refactor them
var profileFieldParsers = []Parser{
    currencyParser,
    modeParser,
    limitParser, 
    categoryParser,
}

func (p *Profile) parseProfile(profilePath string) {
    debug.TraceF("Reading yaml file %v", profilePath)
    file, err := os.Open(profilePath) 
    util.CheckFatalErrorf(err, "Failed to open profile %v", profilePath)
   
    decoder := yaml.NewDecoder(file) 
    data := make(map[interface{}]interface{})
    err = decoder.Decode(data)
    util.CheckFatalError(err)
    debug.TraceF("Parsed yaml: %v", data) 

    for _, profileParser := range profileFieldParsers {
        profileParser.parse(p, data)
    }
}

// Pretty-print current status of profiles
func Dump() {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.Style().Options.SeparateRows = false
    t.Style().Options.SeparateColumns = false
    t.Style().Options.DrawBorder = false

    t.AppendRows([]table.Row{
        {"Current Profile:", currProfile},
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

func (p *Profile) String() string {
    // TODO: maybe put some fmt stuff in display package?
    return fmt.Sprintf("profile:%v\nlimit:%v\nmode:%v\ncategories:%v\n", 
        p.Name,
        p.Limit,
        p.Mode,
        p.Category,
    )
}
