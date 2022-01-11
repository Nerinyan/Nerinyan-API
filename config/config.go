package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	Port         string   `json:"port"`
	TargetDir    string   `json:"targetDir"`
	SlaveServers []string `json:"slaveServers"`
	Sql          struct {
		Id     string `json:"id"`
		Passwd string `json:"passwd"`
		Url    string `json:"url"`
	} `json:"sql"`
	Osu struct {
		Username string `json:"username"`
		Passwd   string `json:"passwd"`
		Token    struct {
			TokenType    string `json:"token_type"`
			ExpiresIn    int64  `json:"expires_in"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"token"`
		BeatmapUpdate struct {
			UpdatedAsc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"updated_asc"`
			UpdatedDesc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"updated_desc"`
			GraveyardAsc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"graveyard_asc"`
		} `json:"beatmapUpdate"`
	} `json:"osu"`
}

var Config config

func LoadConfig() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		out, err := os.Create("./config.json")
		if err != nil {
			panic(err)
		}
		defer out.Close()
		body, err := json.MarshalIndent(Config, "", "    ")
		if err != nil {
			panic(err)
		}
		// Write the body to file
		if _, err = out.Write(body); err != nil {
			panic(err)
		}
	}

	err = json.Unmarshal(b, &Config)
	if err != nil {
		return
	}

}
func (v *config) Save() {
	file, _ := json.MarshalIndent(v, "", "  ")
	_ = ioutil.WriteFile("config.json", file, 0755)
}
