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
  "strings"
)

var shops = []Shop{}

type Shop struct {
  Name string `json:"name"`
  Distance float64 `json:"dist"`
  Y float64
  X float64
}

type OutputShop struct {
  Name string `json:"name"`
  Distance float64 `json:"dist"`
}

func main() {
  http.HandleFunc("/needcoffee", needCoffee)
  http.HandleFunc("/addcoffees", addCoffees)

  http.ListenAndServe(":4000", nil)
}

func init() {
  readShops()
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

func readShops() {
  csvFile, _ := os.Open("data.csv")
  reader := csv.NewReader(bufio.NewReader(csvFile))
  for {
    line, error := reader.Read()
    if error == io.EOF {
        break
    } else if error != nil {
        fmt.Println(error)
    }
    y, _ := strconv.ParseFloat(strings.TrimSpace(line[1]), 64)
    x, _ := strconv.ParseFloat(strings.TrimSpace(line[2]), 64)
    shops = append(shops, Shop{
      Name: line[0],
      Y: y,
      X: x,
    })
  }
}

func needCoffee(w http.ResponseWriter, r *http.Request) {
  xString := r.URL.Query()["x"][0]
  yString := r.URL.Query()["y"][0]

  x, _ := strconv.ParseFloat(xString, 64)
  y, _ := strconv.ParseFloat(yString, 64)

  for i, shop := range shops {
    shops[i].Distance = distance(x, y, shop.X, shop.Y)
  }

  //sort them
  sort.Slice(shops, func(i, j int) bool {
    return shops[i].Distance < shops[j].Distance
  })


  var outputShops []OutputShop
  for _, shop := range shops {
    outputShops = append(outputShops, OutputShop{ Name: shop.Name, Distance: shop.Distance })
  }

  js, err := json.Marshal(outputShops)

  if err != nil {
    fmt.Fprintf(w, "Oh no!", err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func addCoffees(w http.ResponseWriter, r *http.Request) {

  if r.Method != "POST" {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    return
  }

  decoder := json.NewDecoder(r.Body)

  var newShops []Shop
  err := decoder.Decode(&newShops)
  if err != nil {
    fmt.Println(err)
  }

  shops = append(shops, newShops...)

  fmt.Println(shops)
}
