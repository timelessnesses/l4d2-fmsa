# l4d2-fmsa

Left 4 Dead 2 "Fuck Modded Server" Application. A Golang application to help you block those annoying lewd 4 dead servers! [FAQ](#faq)

## Use

Go ahead in Release panel and grab the latest version

## Building
<br>
You need gcc installed and turn `CGO_ENABLED` to 1 then 
`go install` and `go build main.go`  
(Note: this program requires you to install ActiveState' Tk and Tcl. If you are using Windows x64 Operating System those zipped files should be ready for you.)

## Note

this is still a windows only project.

## FAQ
<br>Q: What are those command prompt on openning for first time?<br>
A: its a command prompt for telling firewall to block preset of IPs. Don't worry it will always ask you an administrator permission even with you run the program in administrator mode.<br>
Q: The program crashed! What should I do?<br>
A: Follow this instructions.
1. Clone this repository (or just get the source of the release, there's source code in either tar.gz or zip)
2. Install gcc via scoop or mingw
3. Run `go run .`
4. Try replicate what you just did last time
5. If you did and it errors out, copy entire error stack trace and report it to GitHub and go to issue and report it there


