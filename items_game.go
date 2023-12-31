package main

import (
	"github.com/baldurstod/vdf"
	"encoding/json"
)

type itemMap map[string]*item
type itemGameMap map[string]interface{}

type itemsGame struct {
	medals bool `default:false`
	itemsVDF *vdf.KeyValue
	Prefabs itemMap
	Items itemMap
}

func (this *itemsGame) MarshalJSON() ([]byte, error) {
	ret := make(itemGameMap)

	ret["items"] = *this.getItems()
	//ret["systems"] = *this.MarshalSystems()

	return json.Marshal(ret)
}

/*func (this *itemsGame) MarshalItems() *itemGameMap {
	items := make(itemGameMap)
	for itemId, item := range *this.getItems() {
		styles := item.getStyles()
		//fmt.Println(len(styles))
		if len(styles) > 1 {
			for _, styleId := range styles {
				items[itemId + "~" + styleId] = &itemStyle{it: item, styleId: styleId}
			}
		} else {
			panic("Fix styles")
			//items[itemId] = &itemStyle{it: item, styleId: "0"}
		}
	}

	return &items
}*/

func (this *itemsGame) getItems() (*itemMap) {
	items := make(itemMap)

	for itemId, item := range this.Items {
		items[itemId] = item
	}
	return &items
}

func (this *itemsGame) init(dat []byte) {
	vdf := vdf.VDF{}
	root := vdf.Parse(dat)
	this.itemsVDF, _ = root.Get("items_game")
	this.Prefabs = make(itemMap)
	this.Items = make(itemMap)

	if prefabs, ok := this.itemsVDF.Get("prefabs"); ok {
		for _, val := range prefabs.GetChilds() {
			var it = item{}
			if it.init(this, val) {
				this.Prefabs[it.Id] = &it
			}
		}
	}

	if items, ok := this.itemsVDF.Get("items"); ok {
		for _, val := range items.GetChilds() {
			var it = item{}
			if it.init(this, val) {
				this.Items[it.Id] = &it
			}
		}
	}

	for _, it := range this.Items {
		it.initPrefabs()
	}
}

func (this *itemsGame) getPrefab(prefabName string) *item {
	return this.Prefabs[prefabName]
}

/*func (this *itemsGame) getItemsPerHero() map[string]*hero {
	heroes := make(map[string]*hero)

	for _, item := range this.Items {
		for _, npc := range item.getUsedByHeroes() {
			h, ok := heroes[npc]
			if !ok {
				h = &hero{npc: npc}
				heroes[npc] = h
			}

			h.addItem(item)
		}
	}
	return heroes
}*/
