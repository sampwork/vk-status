package main

import (
    "fmt"
    "gopkg.in/resty.v0"
    "encoding/json"
    "time"
    "flag"
    "github.com/go-ini/ini"
    "github.com/cratonica/trayhost"
	"runtime"
)

var steamStatusName = [7]string {"Offline", "Online", "Busy", "Away", "Snooze", "Trading", "Looking to play"}

type userJSON struct {
    Response struct {
        Players []struct {
            NickName string `json:"personaname"`
            SteamStatus int `json:"personastate"`
            GameName string `json:"gameextrainfo"`
        } `json:"players"`
    } `json:"response"`
}
type steamu struct {
    Nickname string
    Status string
    GameName string
}
func setStatus(vktoken, text string) {

    _, err := resty.R().
      SetQueryParams(map[string]string{
          "access_token": vktoken,
          "text": text,
      }).
      Get("https://api.vk.com/method/status.set?")

    if err != nil {
        fmt.Println(err)
    }
}
 
func getSteamInfo(token, steamid string) (string, string, string) {
    resp, err := resty.R().
      SetQueryParams(map[string]string{
          "key": token,
          "steamids": steamid,
      }).
      Get("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002")

    if err != nil {
        fmt.Println(err)
    }

	var steamuser userJSON    

    if err := json.Unmarshal(resp.Body(), &steamuser); err != nil {
        fmt.Println(err)
    }

    return steamuser.Response.Players[0].NickName, steamuser.Response.Players[0].GameName, steamStatusName[steamuser.Response.Players[0].SteamStatus]
}

func process(steamtoken, steamid, vktoken string) {
	var nickname, gamename, status, text string

	trayhost.SetUrl("https://github.com/sampwork")

	u := new(steamu)

	for range time.Tick(30 * time.Second) {
		nickname, gamename, status = getSteamInfo(steamtoken, steamid)

		if u.Nickname == nickname && u.GameName == gamename && u.Status == status {
			continue
		}
		
		text = "Steam: " + nickname + "(" + status + ")"
		
		if (gamename != "") {
			text += " playing: " + gamename
		}

		setStatus(vktoken, text)

		u.Nickname = nickname
		u.GameName = gamename
		u.Status = status
	}
}

func main() {
	runtime.LockOSThread()

	flag.Parse()

	cfg, err := ini.InsensitiveLoad("settings.ini")

	if err != nil {
		fmt.Println(err)
	}

	steamtoken := cfg.Section("settings").Key("steamtoken").String()
	steamid := cfg.Section("settings").Key("steamid").String()
	vktoken := cfg.Section("settings").Key("vktoken").String()

	if flag.NArg() == 3 {
		steamtoken = flag.Arg(0)
		steamid = flag.Arg(1)
		vktoken = flag.Arg(2)
	} 
	
	if steamtoken == "" || steamid == "" || vktoken == "" {
		fmt.Println("Update settings")
		return
	}

	go process(steamtoken, steamid, vktoken)

	trayhost.EnterLoop("VkStatus", iconData)
}