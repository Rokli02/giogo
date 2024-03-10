## GioGo

Ez egy ___Golang___ nyelven iródott, grafikus asztali alkalmazás, ami a ___Gioui___ grafikus könyvtárat használja fel. </br>
Lényegében egy teszt tesztelgetős, tanulós projektnek indult, de egy aknakereső játék jött ki a végén. </br>
#### Jelenleg támogatott játékmódok:
  - Egyszerú Single player

#### Fejlesztés alatt levő játékmódok:
  - Co-op aknakereső

### Futtatás

  A _Go_ nyelvből adódóan lehetőség van ".exe" fájl létrehozása nélkül futtatni az következő paranccsal, ha a "giogo" mappában vagyunk:
  </br>```go run .```</br>
  Viszont lehetőségünk van egy ".exe" fájl létrehozására is, a következő paranccsal
  </br>```go build -ldflags="-H windowsgui" -o MineGO.exe```</br>