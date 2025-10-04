Require prism
Output "./assets/fails.gif"

Set Width 700
Set Height 800
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
Type "go install ." Enter
Type "clear" Enter
Show

Type "prism -f ./examples/pkg1 ./examples/pkg3" Enter
Sleep 10
