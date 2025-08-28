# OBSST
OBS Simple Toggle

## What? Why?

OBS Simple Toggle is a simple tool that allows you to toggle between two (or more) OBS scenes by creating scene groups. 
The main use is to allow to toggle between two scenes using a single streamdeck button., because weirdly streamdeck doesn't support toggling between two scenes.

## How to use?
1. Create a scene group in config.yaml (see demo file included)
2. List the scenes you want to toggle between in the scene group.
3. Add a button to the streamdeck using the Open action (in the system category)
4. Press the streamdeck button to toggle between the scenes in the scene group.

## How to build?
OBSST is a simple golang application, to build it you need to have a recent version of golang installed.

```
git clone https://github.com/david-ca6/obsst.git
go mod tidy
go build .
```
