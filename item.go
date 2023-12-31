package main

import (
	"strconv"
	"strings"
	"encoding/json"
	"github.com/baldurstod/vdf"
	"fmt"
)

var ITEM_FIELDS = map[string]string{
	"image_inventory": "imageInventory",
	"model_player": "modelPlayer",
	"model_player1": "modelPlayer1",
	"model_player2": "modelPlayer2",
	"model_player3": "modelPlayer3",
	"item_slot": "slot",
	"baseitem": "baseItem",
	"item_rarity": "rarity",
	"workshop_accepted": "workshopAccepted",
	"used_by_classes": "used_by_classes",
}

type item struct {
	ig *itemsGame
	Id string
	Prefab string
	prefabs []*item
	isPrefabsInitialized bool `default:false`
	kv *vdf.KeyValue
}

func (this *item) init(ig *itemsGame, kv *vdf.KeyValue) bool {
	this.ig = ig
	this.Id = kv.Key
	this.kv = kv

	return true
}

func (this *item) initPrefabs() {
	if !this.isPrefabsInitialized {
		this.isPrefabsInitialized = true
		if s, ok := this.kv.GetString("prefab"); ok {
			prefabs := strings.Split(s, " ")
			for _, prefabName := range prefabs {
				prefab := this.ig.getPrefab(prefabName)
				if prefab != nil {
					prefab.initPrefabs() //Ensure prefab is initialized
					this.prefabs = append(this.prefabs, prefab)
				} else {
					fmt.Println("Unknown prefab ", prefabName)
				}
			}
		}
	}
}

func (this *item) getStringAttribute(attributeName string) (string, bool) {
	if s, ok := this.kv.GetString(attributeName); ok {
		return s, true
	}

	for _, prefab := range this.prefabs {
		if s, ok := prefab.getStringAttribute(attributeName); ok && s != "0" {//TODO: remove s != "0"
			return s, true
		}
	}
	return "", false
}

func (this *item) getStringMapAttribute(attributeName string, i *itemStringMap) {
	for _, prefab := range this.prefabs {
		prefab.getStringMapAttribute(attributeName, i)
	}

	if sm, ok := this.kv.GetStringMap(attributeName); ok {
		for key, val := range *sm {
			(*i)[key] = val
		}
	}
}

func (this *item) getUsedByHeroes() []string {
	ret := []string{}

	if usedByHeroes, ok := this.kv.GetStringMap("used_by_heroes"); ok {
		for key, val := range *usedByHeroes {
			if val == "1" {
				ret = append(ret, key)
			}
		}
	}
	return ret
}

func (this *item) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"], _ = strconv.Atoi(this.Id)

	if s, ok := this.getStringAttribute("item_name"); ok {
		ret["name"] = getStringToken(s)
	}

	for inKey, outKey := range ITEM_FIELDS {
		if s, ok := this.getStringAttribute(inKey); ok {
			s = strings.ReplaceAll(s, "\\", "/")
			if inKey == "image_inventory" {
				s = strings.ToLower(s)
			}

			ret[outKey] = s
		} else {
			attribute := make(itemStringMap)
			this.getStringMapAttribute(inKey, &attribute)
			if len(attribute) > 0 {
				ret[outKey] = attribute
			}
		}
	}

	this.MarshalVisuals(&ret)


	if bundle, ok := this.kv.Get("bundle"); ok {
		if sm, ok := bundle.ToStringMap(); ok {
			items := []string{}
			for key, val := range *sm {
				if val == "1" {
					items = append(items, key)
				}
			}
			ret["bundle"] = items
		}
	}

	return json.Marshal(ret)
}

func (this *item) MarshalVisuals(ret *map[string]interface{}) {

	modifiers := []interface{}{}

	if visuals, ok := this.kv.Get("visuals"); ok {
		for _, kv := range visuals.Value.([]*vdf.KeyValue) {
			if strings.HasPrefix(kv.Key, "asset_modifier") {
				modifiers = append(modifiers, kv)
			}
			if strings.HasPrefix(kv.Key, "styles") {
				(*ret)["styles"] = kv
			}
			if strings.HasPrefix(kv.Key, "skin") {
				(*ret)["skin"] = kv
			}
		}
	}

	if len(modifiers) > 0 {
		(*ret)["assetmodifiers"] = modifiers
	}
}
