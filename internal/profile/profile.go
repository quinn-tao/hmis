package profile

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
    "errors"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"gopkg.in/yaml.v3"
)

var (
    ErrNoProfileLoaded = errors.New("No current profile loaded")
)

type Profile struct {
	Name     string
	Limit    coins.RawAmountVal
	Mode     Mode
	Currency coins.Currency
	Category *Category
	profiles []string
	updated  bool
}

// The current loaded profile
var currProfile *Profile

// Profile APIs 
// ===========================================================================

func GetCurrentProfile() (*Profile, error){
    if currProfile == nil {
        return nil, ErrNoProfileLoaded
    }
    return currProfile, nil
}

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
		tokens := strings.Split(file.Name(), ".")
		if len(tokens) < 2 || tokens[1] != "yaml" {
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
	file, err := os.Open(profilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	debug.Tracef("Parsing yaml file %v", profilePath)
	newProfile.ReadFrom(file)

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
		file, err := os.OpenFile(profilePath, os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer file.Close()

		debug.Tracef("Writing back to %v", profilePath)
		err = currProfile.WriteBack(file)
		util.CheckError(err)
	}
	return nil
}

func FindCategory(name string) (c *Category, exists bool) {
	if c, exists := currProfile.Category.FindCategoryWithPath(name); exists {
		return c, exists
	}
	return currProfile.Category.FindCategoryRecursive(name)
}

// Add a new category. This would alter user's profile
func AddCategoryToProfile(path string) error {
	debug.Tracef("Adding category %v", path)
	if _, err := currProfile.Category.AddCategory(path); err != nil {
		return err
	}
	currProfile.updated = true
	return nil
}

// Parsers
// ===========================================================================

// Read and parse a profile from yaml.
//
// Parsing logics are defined via list of field parsers. Each individually
// load a section of users's budgeting profile from generic
// map[interface{}]interface{}
//
// NOTE: order of these parsers matters as some may depends on others.
// For instance, recurrent expenses only make sense after currency settings
// are loaded
func (p *Profile) ReadFrom(file *os.File) error {
	debug.Trace("Constructing yaml map ast from %v", file.Name())

	nameParser := genStringFieldParser(&p.Name, "name")
	var profileFieldParsers = []Parser{
		nameParser,
		currencyParser,
		modeParser,
		limitParser,
		categoryParser,
	}

	decoder := yaml.NewDecoder(file)
	data := make(map[interface{}]interface{})
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	debug.Tracef("Parsed yaml: %v", data)
	for _, profileParser := range profileFieldParsers {
		err = profileParser.parse(p, data)
		if err != nil {
			return err
		}
	}

	return err
}

// Update the curent profile
func (p *Profile) WriteBack(file *os.File) error {
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err := encoder.Encode(p)
	if err != nil {
		return err
	}

	return nil
}

// Printers
// ===========================================================================

func Dump() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Options.SeparateRows = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.DrawBorder = false
	t.SetTitle("Profile Information")

	t.AppendRows([]table.Row{
		{"Name", currProfile.Name},
		{"Mode", currProfile.Mode},
		{"Limit", currProfile.Limit},
		{"Categories", currProfile.Category},
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

func (p Profile) String() string {
	return fmt.Sprintf("profile:%v\nlimit:%v\nmode:%v\ncategories:%v\n",
		p.Name,
		p.Limit,
		p.Mode,
		p.Category,
	)
}
