Output "./assets/colors.gif"

Set Width 1100
Set Height 900
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
Type "go install ." Enter
Type "prism config clear" Enter
Type "clear" Enter
Show

Type "prism --no-color -v ./examples/pkg3" Enter
Sleep 4 

Type "prism --color -f ./examples/pkg1 ./examples/pkg2" Enter
Sleep 4 
