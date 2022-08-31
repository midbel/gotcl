::util::defer {
  puts "leaving mini.tcl"
}

puts $::util::version

proc greet {msg {who foobar}} {
  ::util::defer { puts ">> done greet "}
  puts "$msg $who (from greet)"
}
greet "good morning" vietnam

puts "enter mini.tcl"
set foobar foobar
puts -channel stdout $foobar
puts -channel stderr foo
puts -channel stderr bar
set done "all done ($foobar)"
set mylist [list {fst snd lst}]
puts [::util::typeof $mylist]
unset -nocomplain foobar
llength $mylist
