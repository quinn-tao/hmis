package profile

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/amount"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"gopkg.in/yaml.v3"
)

type Profile struct {
    profiles []string
    Name string
    Limit amount.RawAmountVal
    Mode Mode
    Currency amount.Currency
    Category *Category

    updated bool
    updatedYaml []byte
}

// The current loaded profile 
var currProfile *Profile

// Load current profile 
func LoadProfile() error {
    var newProfile Profile
    dir := config.ProfileDir()
    
    debug.Tracef("Loading profile directory %v", dir)
    profiles, err := os.ReadDir(dir) 
    if err != nil {
        log.Panicf("Error loading profile: %v", err)
        return err
    }
    newProfile.profiles = make([]string, 0)
    
    name := config.CurrentProfileName()
    debug.Tracef("Config: current profile name %v", name)
    for _, file := range profiles {
        tokens := strings.Split(file.Name(),".")
        if len(tokens) < 2 || tokens[1] != "yaml" {
            // TODO: support profile names with multiple dots 
            debug.Tracef("Skipped file %v", file.Name())
            continue
        }

        profile := tokens[0]
        if profile == name {
            newProfile.Name = profile
        }
        newProfile.profiles = append(newProfile.profiles, profile)
    } 

    profilePath := path.Join(dir, newProfile.Name+".yaml")
    newProfile.ReadFrom(profilePath)

    // Read & Parser Yaml
    debug.Trace("Profile config initialized")
    currProfile = &newProfile
    return nil
}

// Close resources related to current profile. Marshals and 
// writes to profile yaml file if user made updates
func UnloadProfile() error {
    if currProfile.updated {
        debug.Tracef("Updating profile %v ...", currProfile.Name)
        dir := config.ProfileDir()
        profilePath := path.Join(dir, currProfile.Name+".yaml")
        err := currProfile.WriteBack(profilePath)
        util.CheckError(err)
    }
    return nil
}

// ****************************************************************************** //
//                                 G E T T E R                                    // 
// ****************************************************************************** //

func FindCategory(name string) (c *Category, exists bool){
    if c, exists := currProfile.Category.FindCategoryWithPath(name); exists {
        return c, exists
    }
    return currProfile.Category.FindCategoryRecursive(name)
}

// ****************************************************************************** //
//                                 S E T T E R                                    // 
// ****************************************************************************** //

// Add a new category. This would alter user's profile
func AddCategoryToProfile(path string) error {
    if _, err := currProfile.Category.AddCategory(path); err != nil {
        return err
    }
    currProfile.updated = true
    return nil 
}

// ****************************************************************************** //
//                                 P A R S E R                                    // 
// ****************************************************************************** //

// List of field parsers. Each individually load a section of 
// users's budgeting profile from generic map[interface{}]interface{}
// NOTE: order of these parsers matters as some may depends on others. 
//       For instance, recurrent expenses only make sense after currency settings 
//       are loaded
var profileFieldParsers = []Parser{
    currencyParser,
    modeParser,
    limitParser, 
    categoryParser,
}

func (p *Profile) ReadFrom(profilePath string) error {
    debug.Tracef("Reading yaml file %v", profilePath)
    file, err := os.Open(profilePath) 
    util.CheckErrorf(err, "Failed to open profile %v", profilePath)
    defer file.Close()
   
    decoder := yaml.NewDecoder(file) 
    data := make(map[interface{}]interface{})
    err = decoder.Decode(data)
    util.CheckError(err)
    debug.Tracef("Parsed yaml: %v", data) 

    for _, profileParser := range profileFieldParsers {
        err = profileParser.parse(p, data)
        if err != nil {
            return err
        }
    }

    return err
}

func (p *Profile) WriteBack(profilePath string) error {
    debug.Tracef("Writing back to %v", profilePath)
    file, err := os.OpenFile(profilePath, os.O_WRONLY | os.O_TRUNC, 0600)
    util.CheckError(err)
    defer file.Close()
    encoder := yaml.NewEncoder(file)
    defer encoder.Close()
    err = encoder.Encode(currProfile)
    if err != nil {
        return err
    }
    return err
}

// ****************************************************************************** //
//                                P R I N T E R                                   //    
// ****************************************************************************** //

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
    
    for _, profile := range currProfile.profiles {
        t.AppendRows([]table.Row{
            {profile},
        })
    }

    t.Render()
}

func (p *Profile) String() string {
    return fmt.Sprintf("profile:%v\nlimit:%v\nmode:%v\ncategories:%v\n", 
        p.Name,
        p.Limit,
        p.Mode,
        p.Category,
    )
}
