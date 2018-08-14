package main

import (
  "net/http"
  "log"
  "html/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" { return }

  index, err := template.ParseFiles("./index.gohtml")
  if err != nil {
    panic(err)
  }

  if err := index.Execute(w, nil); err != nil {
    panic(err)
  }
  //fmt.Fprintln(w, "new home")
}

func main() {
  http.HandleFunc("/", homeHandler)
  log.Fatal(http.ListenAndServe(":8181", nil))
}
