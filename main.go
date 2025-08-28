package main

import (
	"fmt"
	"obsst/obs"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	// fyne
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	Host     string  `mapstructure:"host"`
	Password string  `mapstructure:"password"`
	Port     int     `mapstructure:"port"`
	Groups   []Group `mapstructure:"groups"`
}

type Group struct {
	Name    string   `mapstructure:"name"`
	Type    string   `mapstructure:"type"`
	Scenes  []string `mapstructure:"scenes"`
	Sources []string `mapstructure:"sources"`
}

func main() {
	configPath := ""

	// load env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// handle .app bundle
	exeDir := filepath.Dir(exePath)
	if strings.Contains(exeDir, ".app/Contents/MacOS") {
		appBundlePath := filepath.Dir(filepath.Dir(exeDir))
		configPath = filepath.Dir(appBundlePath) + "/"
	} else {
		configPath = exeDir
	}
	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		a := app.New()
		w := a.NewWindow("obsToggle")
		text := widget.NewLabel("Error reading config file ")
		text1 := widget.NewLabel(configPath)
		w.SetContent(container.NewVBox(text, text1))
		w.Resize(fyne.NewSize(300, 150))
		w.ShowAndRun()
		panic(err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		a := app.New()
		w := a.NewWindow("obsToggle")
		text := widget.NewLabel("Error unmarshalling config file ")
		text1 := widget.NewLabel(configPath)
		w.SetContent(container.NewVBox(text, text1))
		w.Resize(fyne.NewSize(300, 150))
		w.ShowAndRun()
		panic(err)
	}

	group := config.Groups[0].Name
	if len(os.Args) == 2 {
		group = os.Args[1]
	}

	obsToggle(config, group)
}

func obsToggle(config Config, group string) {
	obs, err := obs.New(config.Host, config.Password)
	if err != nil {
		panic(err)
	}

	err = obs.Connect()
	if err != nil {
		panic(err)
	}

	// get current scene
	currentScene, err := obs.GetActiveScene()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current scene: ", currentScene)

	// find group in config
	for _, grp := range config.Groups {
		if grp.Name == group {

			if grp.Type == "scene" {
				for _, scene := range grp.Scenes {
					fmt.Println("Scene: ", scene)
				}

				notInScene := true
				for i, scene := range grp.Scenes {
					if currentScene == scene {
						notInScene = false

						// switch to next scene or first scene
						if scene == grp.Scenes[len(grp.Scenes)-1] {
							err = obs.SetActiveScene(grp.Scenes[0])
							if err != nil {
								panic(err)
							}
						} else {
							err = obs.SetActiveScene(grp.Scenes[i+1])
							if err != nil {
								panic(err)
							}
						}
					}
				}

				if notInScene {
					err = obs.SetActiveScene(grp.Scenes[0])
					if err != nil {
						panic(err)
					}
				}

			} else if grp.Type == "source" {
				for _, source := range grp.Sources {
					fmt.Println("Source: ", source)
				}

				if currentScene != grp.Scenes[0] {
					err = obs.SetActiveScene(grp.Scenes[0])
					if err != nil {
						panic(err)
					}

					err = obs.SetSceneItemEnabled(grp.Scenes[0], grp.Sources[0], true)
					if err != nil {
						panic(err)
					}
				} else {

				}
			}
		}
	}
}
