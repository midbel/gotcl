defer { puts "leaving mini.tcl" }
puts "enter mini.tcl"
set foobar foobar
puts $foobar
puts foo
puts bar
set done "all done ($foobar)"
set mylist [list {fst snd lst}]
puts [typeof $mylist]
llength $mylist
