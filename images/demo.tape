Output images/demo.gif
Set FontSize 18
Set Width 1400
Set Height 900
Set Theme "Builtin Solarized Dark"
Set TypingSpeed 40ms
Set Shell "bash"

Type "cd /home/s3/go-jmp.nigelhamilton.com"
Enter
Sleep 500ms

# jmp in — search in files for keywords
Type "jmp in 'FindFiles'"
Enter
Sleep 4s
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Right
Sleep 3s
Type "q"
Sleep 1s

# jmp on — jump on files in command output
Type "jmp on 'tree -f'"
Enter
Sleep 4s
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Right
Sleep 3s
Type "q"
Sleep 1s

# jmp to — locate a file on the filesystem
Type "jmp to main.go"
Enter
Sleep 4s
Down
Sleep 500ms
Right
Sleep 3s
Type "q"
Sleep 1s
