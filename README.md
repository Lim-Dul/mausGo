# mausGo
Character Generator for Mausritter written in Go.
Used it to learn the language. :)
## Usage
``` go run .\mausGo.go ```
or
``` mausGo.exe ```
after building
### Advanced Options
Use the following flags to set minimum values for certain attributes:
```
  -minSTR int
        Minimum STR (default 2)
  -minDEX int
        Minimum DEX (default 2)
  -minWIL int
        Minimum WIL (default 2)
  -minHP int
        Minimum HP (default 1)
  -minPIPS int
        Minimum Pips (default 1)
  -h
        Print help
```
Example:
```
mausGo.exe -minSTR 12 -minHP 2
```
