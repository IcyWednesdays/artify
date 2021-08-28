package main

type ZoneBoundary struct {
	Start int64
	End   int64
}

type ZoneConfig struct {
	Primary   []ZoneBoundary
	Secondary []ZoneBoundary
	Tertiary  []ZoneBoundary
}

type AvailableScenes map[string]*ZoneConfig

func getSceneByName(sceneName string) *ZoneConfig {
	availableScenes := AvailableScenes{
		"INSIDETOOUT": {
			Primary:   []ZoneBoundary{{Start: 8, End: 15}},
			Secondary: []ZoneBoundary{{Start: 4, End: 7}},
			Tertiary:  []ZoneBoundary{{Start: 0, End: 3}},
		},
		"BLENDED": {
			Primary:   []ZoneBoundary{{Start: 0, End: 2}, {Start: 9, End: 11}},
			Secondary: []ZoneBoundary{{Start: 3, End: 5}, {Start: 12, End: 14}},
			Tertiary:  []ZoneBoundary{{Start: 6, End: 8}, {Start: 15, End: 15}},
		},
	}

	if _, ok := availableScenes[sceneName]; ok {
		return availableScenes[sceneName]
	}
	return availableScenes["INSIDETOOUT"]
}
