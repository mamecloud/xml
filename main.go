package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"encoding/xml"
)

type Mame struct {
	XMLName		xml.Name	`xml:"mame"`
	Build		string		`xml:"build,attr"`
	Machine		[]Machine	`xml:"machine"`
}

type Machine struct {
	XMLName			xml.Name	`xml:"machine"`
	Name			string		`xml:"name,attr"`
	IsMechanical	string		`xml:"ismechanical,attr"`
	IsBios			string		`xml:"isbios,attr"`
	IsDevice		string		`xml:"isdevice,attr"`
	Runnable		string		`xml:"runnable,attr"`
	Description		string		`xml:"description"`
	Year			string		`xml:"year"`
	Rom				[]Rom		`xml:"rom"`
	Driver			Driver		`xml:"driver"`
	Feature			[]Feature	`xml:"feature"`
}

type Rom struct {
	XMLName		xml.Name	`xml:"rom"`
	Name		string		`xml:"name,attr"`
	Size		string		`xml:"size,attr"`
	Crc			string		`xml:"crc,attr"`
	Sha1		string 		`xml:"sha1,attr"`
}

type Driver struct {
	XMLName		xml.Name	`xml:"driver"`
	Status		string		`xml:"status,attr"`		// good, imperfect, preliminary
	Emulation	string		`xml:"emulation,attr"`	// good, preliminary
	Cocktail	string		`xml:"cocktail,attr"`	// preliminary
	Savestate	string		`xml:"savestate,attr"`	// supported, unsupported
}

type Feature struct {
	XMLName		xml.Name	`xml:"feature"`
	Type		string		`xml:"type,attr"`		// camera:54 controls:16 disk:49 graphics:2931 keyboard:13 lan:205 microphone:5 mouse:1 palette:286 printer:68 protection:275 sound:5774 tape:1 timing:330
	Status		string		`xml:"status,attr"`		// imperfect, unemulated
	Overall		string		`xml:"overall,attr"`	// imperfect, unemulated
}

func main() {

    // Open our xmlFile
    xmlFile, err := os.Open("mame.xml")
    // if we os.Open returns an error then handle it
    if err != nil {
		log.Fatal(err)
    }
    defer xmlFile.Close()

	fmt.Println("Successfully Opened mame.xml")
	
	// read our opened xmlFile as a byte array.
	// byteValue, _ := ioutil.ReadAll(xmlFile)

	// var mame Mame
	// xml.Unmarshal(byteValue, &mame)

	// for i, machine := range mame.Machine {
	// 	fmt.Printf("%d: %s", i, machine.Description)
	// }

	// fmt.Println(mame.Build)
	// fmt.Println(len(mame.Machine))

  
	d := xml.NewDecoder(xmlFile)
	total := 0
	count := 0
	mechanical := 0
	device := 0
	bios := 0
	driverStatus := make(map[string]int)
	driverStatusGood := 0
	driverStatusImperfect := 0
	driverStatusPreliminary := 0
	driverEmulation := make(map[string]int)
	driverCocktail := make(map[string]int)
	driverSavestate := make(map[string]int)
	featureType := make(map[string]int)
	featureStatus := make(map[string]int)
	featureOverall := make(map[string]int)
	roms := make(map[string]int)
	romSize := make(map[string]string)
	// romCrc := make(map[string]string)
	// romSha1 := make(map[string]string)
	sizeClash := make([]string, 0)
	// crcClash := make([]string, 0)
	// sha1Clash := make([]string, 0)
	for {
	  tok, err := d.Token()
	  if tok == nil || err == io.EOF {
		// EOF means we're done.
		break
	  } else if err != nil {
		log.Fatalf("Error decoding token: %s", err)
	  }
  
	  switch ty := tok.(type) {
	  case xml.StartElement:
		if ty.Name.Local == "machine" {
			total++
			// If this is a start element named "machine", parse this element
			// fully.
			var machine Machine
			if err = d.DecodeElement(&machine, &ty); err != nil {
			log.Fatalf("Error decoding item: %s", err)
			}

			if machine.IsMechanical == "yes" {
				mechanical++
			} else if machine.IsBios == "yes" {
				bios++
			} else if machine.IsDevice == "yes" {
				device++
			} else {
				count++
				if count % 1000 == 0 {
					fmt.Printf("%s: %s\n", machine.Name, machine.Description)
					for _, rom := range(machine.Rom) {
						fmt.Printf(" - %s %s/%s/%s\n", rom.Name, rom.Size, rom.Crc, rom.Sha1)
					}
				}

				driver := machine.Driver
				driverStatus[driver.Status] = driverStatus[driver.Status] + 1
				driverEmulation[driver.Emulation] = driverEmulation[driver.Emulation] + 1
				driverCocktail[driver.Cocktail] = driverCocktail[driver.Cocktail] + 1
				driverSavestate[driver.Savestate] = driverSavestate[driver.Savestate] + 1
				if driver.Status == "good" {
					driverStatusGood++
				} else if driver.Status == "imperfect" {
					driverStatusImperfect++
				} else if driver.Status == "preliminary" {
					driverStatusPreliminary++
				}
				for _, feature := range(machine.Feature) {
					featureType[feature.Type] = featureType[feature.Type] + 1
					featureStatus[feature.Status] = featureStatus[feature.Status] + 1
					featureOverall[feature.Overall] = featureOverall[feature.Overall] + 1
				}

				for _, rom := range machine.Rom {
					name := rom.Name
					id := fmt.Sprintf("%s-%s-%s-%s", rom.Name, rom.Size, rom.Crc, rom.Sha1)
					roms[name] = roms[name] + 1
					if romSize[name] == "" {
						romSize[name] = rom.Size
					} else if romSize[name] != rom.Size {
						sizeClash = append(sizeClash, id)
					}
				}
			}
		}
	  default:
	  }
	}
  
	fmt.Printf("driverStatus: %v\n", driverStatus)
	fmt.Printf("driverEmulation: %v\n", driverEmulation)
	fmt.Printf("driverCocktail: %v\n", driverCocktail)
	fmt.Printf("driverSavestate: %v\n", driverSavestate)
	fmt.Printf("featureType: %v\n", featureType)
	fmt.Printf("featureStatus: %v\n", featureStatus)
	fmt.Printf("featureOverall: %v\n", featureOverall)
	fmt.Println("count =", count)
	fmt.Println("mechanical =", mechanical)
	fmt.Println("bios =", bios)
	fmt.Println("device =", device)
	fmt.Println("total =", total)
	fmt.Println("good =", driverStatusGood)
	fmt.Println("imperfect =", driverStatusImperfect)
	fmt.Println("preliminary =", driverStatusPreliminary)
	for id, count := range(roms) {
		if count > 1 {
			fmt.Printf("key[%s] value[%d]\n", id, count)
		}
	}
	fmt.Printf("Total multi-roms: %d\n", len(roms))
	if len(sizeClash) > 0 {
		fmt.Println(sizeClash)
	}
  }

