package config

import "github.com/spf13/viper"

func ProfileDir() string {
    return viper.GetString("profile.dir")
}

func CurrentProfileName() string {
    return viper.GetString("profile.name")
}

func StorageLocation() string {
    return viper.GetString("storage_location")
}
