# mausGo
Character Generator for Mausritter written in Go
## Usage
``` go run .\mausGo.go ```
or
``` mausGo.exe ```
after building
### Advanced Options
Use the following flags to set minimum values for certain attributes:
```
-minDEX int
        Minimum DEX (default 2)
  -minHP int
        Minimum HP (default 1)
  -minPIPS int
        Minimum Pips (default 1)
  -minSTR int
        Minimum STR (default 2)
  -minWIL int
        Minimum WIL (default 2)
```
Example:
``` mausGo.exe -minSTR 12 -minHP 2 ```
