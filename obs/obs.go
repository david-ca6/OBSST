package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

type Obs struct {
	client   *goobs.Client
	host     string
	password string

	ObsVersion string
}

func New(host string, password string) (*Obs, error) {
	var err error

	obs := Obs{
		client:   nil,
		host:     host,
		password: password,
	}

	// Connect to OBS to test if the connection is working
	err = obs.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OBS: %w", err)
	}

	defer obs.Disconnect()

	version, err := obs.client.General.GetVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get OBS version: %w", err)
	}

	fmt.Printf("OBS Studio version: %s\n", version.ObsVersion)
	fmt.Printf("Server protocol version: %s\n", version.ObsWebSocketVersion)
	fmt.Printf("Client protocol version: %s\n", goobs.ProtocolVersion)
	fmt.Printf("Client library version: %s\n", goobs.LibraryVersion)

	obs.ObsVersion = version.ObsVersion

	return &obs, nil
}

// *****************************************************************************
// Connection
// *****************************************************************************

func (o *Obs) Connect() error {
	var err error
	o.client, err = goobs.New(o.host, goobs.WithPassword(o.password))
	if err != nil {
		return fmt.Errorf("failed to connect to OBS: %w", err)
	}

	return nil
}

func (o *Obs) Disconnect() {
	o.client.Disconnect()
}

// *****************************************************************************
// Stream
// *****************************************************************************

func (o *Obs) StartStream() error {
	_, err := o.client.Stream.StartStream(nil)
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}
	return nil
}

func (o *Obs) StopStream() error {
	_, err := o.client.Stream.StopStream(nil)
	if err != nil {
		return fmt.Errorf("failed to stop stream: %w", err)
	}
	return nil
}

func (o *Obs) ToggleStream() error {
	_, err := o.client.Stream.ToggleStream(nil)
	if err != nil {
		return fmt.Errorf("failed to toggle stream: %w", err)
	}
	return nil
}

// *****************************************************************************
// Scenes
// *****************************************************************************

func (o *Obs) GetSceneList() ([]string, error) {
	scenesList, err := o.client.Scenes.GetSceneList()
	if err != nil {
		return nil, fmt.Errorf("failed to get scene list: %w", err)
	}
	var sceneNames []string
	for _, scene := range scenesList.Scenes {
		sceneNames = append(sceneNames, scene.SceneName)
	}

	for i, j := 0, len(sceneNames)-1; i < j; i, j = i+1, j-1 {
		sceneNames[i], sceneNames[j] = sceneNames[j], sceneNames[i]
	}

	return sceneNames, nil
}

func (o *Obs) GetActiveScene() (string, error) {
	activeScene, err := o.client.Scenes.GetCurrentProgramScene()
	if err != nil {
		return "", fmt.Errorf("failed to get active scene: %w", err)
	}
	return activeScene.SceneName, nil
}

func (o *Obs) SetActiveScene(sceneName string) error {
	_, err := o.client.Scenes.SetCurrentProgramScene(
		&scenes.SetCurrentProgramSceneParams{
			SceneName: &sceneName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to set active scene: %w", err)
	}
	return nil
}

// *****************************************************************************
// Scene Items
// *****************************************************************************

func (o *Obs) GetSceneItemList(scene string) ([]string, error) {
	params := sceneitems.NewGetSceneItemListParams().WithSceneName(scene)
	items, err := o.client.SceneItems.GetSceneItemList(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get scene item list: %w", err)
	}
	var sceneItemList []string
	for _, v := range items.SceneItems {
		sceneItemList = append(sceneItemList, v.SourceName)
	}

	for i, j := 0, len(sceneItemList)-1; i < j; i, j = i+1, j-1 {
		sceneItemList[i], sceneItemList[j] = sceneItemList[j], sceneItemList[i]
	}

	return sceneItemList, nil
}

func (o *Obs) GetSceneItemID(sceneName string, itemName string) (int, error) {
	params := sceneitems.NewGetSceneItemListParams().WithSceneName(sceneName)
	items, err := o.client.SceneItems.GetSceneItemList(params)
	if err != nil {
		return 0, fmt.Errorf("failed to get scene item list: %w", err)
	}
	for _, v := range items.SceneItems {
		if v.SourceName == itemName {
			return v.SceneItemID, nil
		}
	}
	return 0, nil
}

func (o *Obs) SetSceneItemEnabled(sceneName string, sceneItem string, enabled bool) error {
	itemID, err := o.GetSceneItemID(sceneName, sceneItem)
	if err != nil {
		return fmt.Errorf("failed to get scene item id: %w", err)
	}

	_, err = o.client.SceneItems.SetSceneItemEnabled(
		&sceneitems.SetSceneItemEnabledParams{
			SceneName:        &sceneName,
			SceneItemId:      &itemID,
			SceneItemEnabled: &enabled,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to set scene item enabled: %w", err)
	}
	return nil
}
