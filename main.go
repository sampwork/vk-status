package main

import (
    "fmt"
    "gopkg.in/resty.v0"
    "encoding/json"
    "time"
    "github.com/go-ini/ini"
)

var steamStatusName = [7]string {"Offline", "Online", "Busy", "Away", "Snooze", "Trading", "Looking to play"}

type steamUser struct {
    nickName string
    gameName string
    statusID int
}

type iniFile struct {
    tokenSteam  string
    tokenVk     string
    idSteam64   string
}

func (i *iniFile) Load(filename string) error {
    file, err := ini.InsensitiveLoad(filename)

    if err != nil {
        return err
    }

    i.tokenVk = file.Section("settings").Key("vktoken").String()
    i.tokenSteam = file.Section("settings").Key("steamtoken").String()
    i.idSteam64 = file.Section("settings").Key("steamid").String()

    return nil
}

func (steam *steamUser) updateData(token, id string) (bool, bool, error){
    resp, err := resty.R().SetQueryParams(map[string]string{"key": token, "steamids": id}).Get("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?")

    if err != nil {
        return false, false, err
    }

    if resp.StatusCode() != 200 {
        fmt.Println("steam status code: ", resp.StatusCode())
        return false, true, nil
    }
    type userJSON struct {
        Response struct {
            Players []struct {
                Name string `json:"personaname"`
                Status int `json:"personastate"`
                GameName string `json:"gameextrainfo"`
            } `json:"players"`
        } `json:"response"`
    }

	var steamuser userJSON    

    if err := json.Unmarshal(resp.Body(), &steamuser); err != nil {
        return false, false, err
    }    

    if  steam.nickName == steamuser.Response.Players[0].Name && 
        steam.gameName == steamuser.Response.Players[0].GameName && 
        steam.statusID == steamuser.Response.Players[0].Status {
        return true, false, nil
    }

    steam.nickName = steamuser.Response.Players[0].Name
    steam.gameName = steamuser.Response.Players[0].GameName
    steam.statusID = steamuser.Response.Players[0].Status

    return false, false, nil
}

func main() {
    fmt.Println("started...")

    ini := new(iniFile)

    err := ini.Load("settings.ini")

    if err != nil {
        fmt.Println("inifile: ", err)
    }

    steam := new(steamUser)

    for range time.Tick(30 * time.Second) {
        doNotUpdate, status, err := steam.updateData(ini.tokenSteam, ini.idSteam64)

        if err != nil {
            fmt.Println("get steam data err: ", err)
            continue
        }

        if doNotUpdate || status {
            continue
        }

        statusText := "Steam: " + steam.nickName + "(" + steamStatusName[steam.statusID] + ")"

        if steam.gameName != "" {
            statusText += " playing: " + steam.gameName
        }

        _, err = resty.R().SetQueryParams(map[string]string{"access_token": ini.tokenVk, "text": statusText}).Get("https://api.vk.com/method/status.set?")

        if err != nil {
            fmt.Println("set vk status err: ", err)
        }
    }
}
