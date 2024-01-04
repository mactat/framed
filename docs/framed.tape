Output ./docs/static/demo.gif

# Settings

Set FontSize 20
Set Width 1200
Set Height 600

# Requirements

Require framed

# Setup

Hide
Type "mkdir -p /tmp/framed-demo"
Enter
Type "cd /tmp/framed-demo"
Enter
Ctrl+L
Show

# Flow

Sleep 1s
Type "echo 'Welcome to Framed!'"
Enter
Sleep 4s
Ctrl+L

Sleep 1s
Type "tree"
Enter
Sleep 4s
Ctrl+L

Sleep 1s
Type "framed import --example python"
Enter
Sleep 4s
Ctrl+L

Sleep 1s
Type "tree"
Enter
Sleep 4s
Ctrl+L

Sleep 1s
Type "cat framed.yaml"
Enter
Sleep 6s
Ctrl+L

Sleep 1s
Type "framed verify"
Enter
Sleep 6s
Ctrl+L

Sleep 1s
Type "framed create --files"
Enter
Sleep 6s
Ctrl+L

Sleep 1s
Type "framed verify"
Enter
Sleep 4s
Ctrl+L

Sleep 1s
Type "tree"
Enter
Sleep 5s
Ctrl+L

Sleep 1s
Type "rm setup.py"
Enter
Sleep 1s
Type "framed verify"
Enter
Sleep 6s
Ctrl+L

# Clean
Hide
Type "cd -"
Enter
Type "rm -r /tmp/framed-demo"
Enter