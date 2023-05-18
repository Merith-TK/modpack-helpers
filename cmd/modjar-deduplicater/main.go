package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Fabric/Quilt Mod Json
type FabricModJson struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}
type QuiltModJson struct {
	QuiltLoader struct {
		ID       string `json:"id"`
		Version  string `json:"version"`
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"quilt_loader"`
}
type ForgeModToml struct {
	Mods []struct {
		ModID       string `toml:"modId"`
		Version     string `toml:"version"`
		DisplayName string `toml:"displayName"`
	} `toml:"mods"`
}

var (
	fabricMods []FabricModJson
	quiltMods  []QuiltModJson
	forgeMods  []ForgeModToml

	cleanPattern = regexp.MustCompile(`[<>;:\"|?*]`)
)

// deduplicate mod jars based off the name of the jar
// if the jar name is the same, then we will only keep the one with the highest version
// or the latest modified date
func main() {
	// get list of jars in ./mods
	files, err := ioutil.ReadDir("./mods")
	if err != nil {
		log.Fatalln("[READ DIR]", err)
	}

	// for modjars,
	for _, file := range files {
		// filter out non-jars
		if !strings.HasSuffix(file.Name(), ".jar") {
			log.Println("[SKIP]", file.Name())
			continue
		}

		// read the mod.json file and return the name, version, and filename
		name, version, loader := readMod("./mods/" + file.Name())
		name = cleanUnicode(name)
		name = cleanPattern.ReplaceAllString(name, "")

		if loader != "" {
			newname := name + "-" + version + ".jar"
			if newname == file.Name() {
				continue
			}
			err := os.Rename("./mods/"+file.Name(), "./mods/"+newname)
			if err != nil {
				log.Fatalln("[RENAME MOD]", err)
			} else {
				log.Println("[RENAME MOD]", file.Name(), ">>", newname)
			}
		} else {
			log.Println("Unsupported Loader", file.Name())
		}
	}
	for _, mod := range fabricMods {
		// check if there is a duplicate in the list
		for _, mod2 := range fabricMods {
			if mod.Name == mod2.Name && mod.Version != mod2.Version {
				// if there is a duplicate, delete the one with the lower semver version
				if mod.Version < mod2.Version {
					name := cleanPattern.ReplaceAllString(mod.Name, "")
					err := os.Remove("./mods/" + name + "-" + mod.Version + ".jar")
					if err != nil {
						log.Fatalln("[REMOVE MOD]", err)
					}
				}
			}
		}
	}
}

// read the mod.json file and return the name, version, and loader
func readMod(filename string) (string, string, string) {
	var fabricJson FabricModJson
	var quiltJson QuiltModJson
	var forgeMod ForgeModToml
	var loader string
	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatalln("[READ ZIP]", err)
	}
	defer zipReader.Close()
	for _, file := range zipReader.File {
		contents, _ := file.Open()

		// * Clean stray \n characters from the json
		var buf bytes.Buffer
		_, err := buf.ReadFrom(contents)
		if err != nil {
			log.Fatalln("[READ FROM]", err)
		}
		cleanedContents := strings.ReplaceAll(buf.String(), "\n", "")
		// * End Clean

		if file.Name == "fabric.mod.json" {
			loader = "fabric"
			fabricJsonErr := json.Unmarshal([]byte(cleanedContents), &fabricJson)
			if fabricJsonErr != nil {
				log.Println("[DECODE ERROR]", filename)
				log.Fatalln("[DECODE FABRIC MOD]", fabricJsonErr)
			}
			break
		}
		if file.Name == "quilt.mod.json" {
			loader = "quilt"
			quiltJsonErr := json.Unmarshal([]byte(cleanedContents), &quiltJson)
			if quiltJsonErr != nil {
				log.Println(filename)
				log.Println("[DECODE ERROR]", filename)
				log.Fatalln("[DECODE QUILT MOD]", quiltJsonErr)
			}
			break
		}
		if file.Name == "META-INF/mods.toml" {
			loader = "forge"
			forgeTomlErr := toml.NewDecoder(contents).Decode(&forgeMod)
			if forgeTomlErr != nil {
				log.Println("[DECODE ERROR]", filename)
				log.Fatalln("[DECODE FORGE MOD]", forgeTomlErr)
			}
			// NOTICE: We don't break here because we want to check if there is a fabric/quilt mod.json,
			// if there is, then we will use that instead of the forge mod.toml due to inconsistencies
			// between the forge mod.tomls
		}
	}

	if loader == "fabric" {
		fabricMods = append(fabricMods, fabricJson)
		return fabricJson.Name, fabricJson.Version, loader
	}
	if loader == "quilt" {
		quiltMods = append(quiltMods, quiltJson)
		return quiltJson.QuiltLoader.Metadata.Name, quiltJson.QuiltLoader.Version, loader
	}
	if loader == "forge" {
		if len(forgeMod.Mods) != 0 {
			forgeMods = append(forgeMods, forgeMod)
			log.Println("Forge Mod", filename)
			return forgeMod.Mods[0].DisplayName, forgeMod.Mods[0].Version, loader
		}
	}
	return "", "", ""
}
