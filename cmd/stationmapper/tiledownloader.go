package main

import (
  "errors"
  "fmt"
  "io"
  "math"
  "net/http"
  "os"
  "time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

var tileServer = "tile.openstreetmap.org"
var rootpath = "/var/data/stationmapper/osm"

func mkdir(d string){

  err := os.MkdirAll(d, 0755)
  check(err)
}

func downloadTilesForZoom(zoom int){
  tiles := math.Pow(4, float64(zoom))
  axis := int(math.Sqrt(tiles))

  fmt.Printf("Downloading tiles for zoom [%d], %d tiles\n", zoom, int(tiles))
  
  for x := 0; x < axis; x++ {
    for y := 0; y < axis; y++ {
      res, err := tileExists(zoom, x, y)
      check(err)

      if !res {
        err := downloadTile(zoom, x, y)
        check(err)
        time.Sleep(time.Duration(zoom) * time.Second)
      }
    }
  }
}

func tileExists(zoom int, x int, y int) (bool, error) {

  filepath := fmt.Sprintf("%s/%d/%d/%d.png",
    rootpath, zoom, x, y)
    
  if _, err := os.Stat(filepath); err == nil {
    return true, nil

  } else if errors.Is(err, os.ErrNotExist) {
    // path/to/whatever does *not* exist
    return false, nil

  } else {
    return false, err
  }

}


func downloadTile(zoom int, x int, y int) error {

  filepath := fmt.Sprintf("%s/%d/%d/%d.png",
    rootpath, zoom, x, y)

  // Get the file
  url := fmt.Sprintf("http://%s/%d/%d/%d.png",
    tileServer, zoom, x, y)

  fmt.Printf("Downloading [%s]\n", url)

  client := &http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  if err != nil { return err }

  req.Header.Set("User-Agent", "Golang/StationMapper github.com/mouckatron/portableops-stationmapper")
  
  resp, err := client.Do(req)
  if err != nil { return err }
  defer resp.Body.Close()

  // Create the directory
  mkdir(fmt.Sprintf("%s/%d/%d", rootpath, zoom, x))
  
  // Create the file
  out, err := os.Create(filepath)
  if err != nil {
    return err
  }
  defer out.Close()

  // Write the file
  _, err = io.Copy(out, resp.Body)
  
  return err
}

// GetTiles downloads a set of tiles from $tileserver to $rootpath
func GetTiles() {
  
  mkdir(rootpath)
  
  for z := 0; z <= 5 ; z++ {
    downloadTilesForZoom(z)
  }
}
