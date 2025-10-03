Require prism
Output "./assets/fails.gif"

Set Width 1000
Set Height 1200
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
Type "go build ." Enter
Type "go install ." Enter
Type "clear" Enter
Show

Type "prism -f ./examples/pkg1 ./examples/pkg3" Enter
Sleep 5
