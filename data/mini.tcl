::util::defer {
  puts "leaving mini.tcl"
}

puts $::util::version

proc greet {} {
  ::util::defer { puts ">> done greet "}
  puts "hello foobar (from greet)"
}
greet

puts "enter mini.tcl"
set foobar foobar
puts $foobar
puts foo
puts bar
set done "all done ($foobar)"
set mylist [list {fst snd lst}]
puts [::util::typeof $mylist]
llength $mylist
