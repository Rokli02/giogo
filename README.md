## GioGo
  Ez egy ___Golang___ nyelven iródott, grafikus asztali alkalmazás, ami a ___Gioui___ grafikus könyvtárat használja fel. </br>
  Lényegében egy tesztelgetős, tanulós projektnek indult, de egy aknakereső játék jött ki a végén. </br>

#### Jelenleg támogatott játékmódok:
  - `Egyszerú Single player`
  - `Co-op aknakereső`

#### Fejlesztés alatt levő játékmódok:
  - `-`

### Futtatás
  A _Go_ nyelvből adódóan lehetőség van ".exe" fájl létrehozása nélkül futtatni az következő paranccsal, ha a "giogo" mappában vagyunk:
  </br>```go run .```</br>
  Viszont lehetőségünk van egy ".exe" fájl létrehozására is, a következő paranccsal
  </br>```go build -ldflags="-H windowsgui -s -w" -o MineGO.exe```</br>

### Paraméterek
  Lehetőség van induláskor paraméterek megadására, amivel tudjuk állítani milyen __nehézségü__ legyen a játék,
  valamint **név beállítás**ára, amit Lobby-ba csatlakozáskor láthatunk

  -u, --username: Felhasználónév beállítása
  -w, --width:    Egyjátékos mód pályának a szélessége
  -h, --height:   Egyjátékos mód pályának a magassága
  -m, --mines:    Egyjátékos mód pályának az aknák száma
