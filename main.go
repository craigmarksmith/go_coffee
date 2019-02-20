package main

import (
  "io"
  "bufio"
  "fmt"
  "net/http"
  "encoding/csv"
  "encoding/json"
  "math"
  "strconv"
  "sort"
  "os"
)

type Shop struct {
  Name string
  Y float64
  X float64
}

type ShopDistance struct {
  Name string
  Distance float64
}

func main() {
  http.HandleFunc("/needcoffee", needCoffee)
  http.HandleFunc("/addcoffees", addCoffees)

  http.ListenAndServe(":4000", nil)
}

func distance(x1, y1, x2, y2 float64) float64 {
  newX := math.Pow((x1-x2), 2)
  newY := math.Pow((y1-y2), 2)
  sum := newX+newY
  sumSq := math.Sqrt(sum)
  return (toFixed(sumSq, 4))
}

func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}

func readShops() []Shop {
  var shops = []Shop{}

  csvFile, _ := os.Open("data.csv")
  reader := csv.NewReader(bufio.NewReader(csvFile))
  for {
    line, error := reader.Read()
    if error == io.EOF {
        break
    } else if error != nil {
        fmt.Println(error)
    }
    y, _ := strconv.ParseFloat(line[1], 64)
    x, _ := strconv.ParseFloat(line[2], 64)
    shops = append(shops, Shop{
      Name: line[0],
      Y: y,
      X: x,
    })
  }

  return shops
}

func needCoffee(w http.ResponseWriter, r *http.Request) {

  shops := readShops()

  xString := r.URL.Query()["x"][0]
  yString := r.URL.Query()["y"][0]

  x, _ := strconv.ParseFloat(xString, 64)
  y, _ := strconv.ParseFloat(yString, 64)

  //array of distances
  shopDistances := []ShopDistance{}
  for i := 0; i < len(shops); i++ {
    shopDistance := ShopDistance{shops[i].Name, distance(x, y, shops[i].X, shops[i].Y)}
    shopDistances = append(shopDistances, shopDistance)
  }  

  //sort them
  sort.Slice(shopDistances, func(i, j int) bool {
    return shopDistances[i].Distance < shopDistances[j].Distance
  })

  js, err := json.Marshal(shopDistances[:3])

  if err != nil {
    fmt.Fprintf(w, "Oh no!", err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func addCoffees(w http.ResponseWriter, r *http.Request) {

  if r.Method == "POST" {
    decoder := json.NewDecoder(r.Body)

    var newShops []Shop
    err := decoder.Decode(&newShops)

    file, ferr := os.OpenFile("data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if ferr != nil {
      fmt.Println(ferr)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for i := 0; i < len(newShops); i++ {
      newLine := []string{newShops[i].Name, fmt.Sprintf("%f", newShops[i].Y), fmt.Sprintf("%f", newShops[i].X)}
      writer.Write(newLine)
    }

    if err != nil {
      fmt.Println(err)
    }
  } else {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
  }

}