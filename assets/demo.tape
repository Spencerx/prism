Require prism
Output "./assets/demo.gif"

Set Width 1000
Set Height 1200
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
Type "go build ." Enter
Type "go install ." Enter
Type "clear" Enter
Show

Type "prism" Enter
Sleep 5

Type "prism -v ./examples/pkg1" Enter
Sleep 5

Type "prism -f ./examples/pkg1 ./examples/pkg3" Enter
Sleep 5

Type "prism ./examples/empty_dir" Enter
Sleep 5
