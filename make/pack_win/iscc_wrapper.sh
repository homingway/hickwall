#!/bin/sh  
unset DISPLAY  
scriptname=$1  
[ -f "$scriptname" ] && scriptname=$(winepath -w "$scriptname")  
wine "C:\Program Files (x86)\Inno Setup 5\ISCC.exe" "$scriptname" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9"