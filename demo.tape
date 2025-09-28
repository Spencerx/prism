Require prism
Output "demo.gif"

Set Width 1200
Set Height 1400
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
Type "go build ." Enter
Type "go install ." Enter
Type "clear" Enter
Show

Type "prism" Enter
Sleep 5

Type "prism -v ./pkg3" Enter
Sleep 5

Type "prism -f ./pkg1 ./pkg3" Enter
Sleep 5
