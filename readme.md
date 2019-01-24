```
$ extdir -h
extdir recurses a directory and reports file counts by extension.
Version 0.1
Usage: extdir PATH
PATH is the directory to start the recurse.
  -path
        Print full paths
  -version
        Print version and exit

$ extdir .
File count: 4

         # extension
         1 .
         1 .exe
         1 .go
         1 .md
         
```