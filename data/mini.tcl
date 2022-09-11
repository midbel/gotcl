::util::defer {
  puts "leaving mini.tcl"
}

echo "echo from shell!!!"

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

puts [::util::help puts]

set done "all done ($foobar)"
set mylist [list {fst snd lst}]
puts [string tolower "HELLO WORLD"]
puts [string toupper "hello world"]
puts [::util::typeof $mylist]
unset -nocomplain foobar
llength $mylist

set count 1
puts "count is: $count"
proc printCount {} {
  upvar count local
  puts "printCount: $local"
  incr local
}
printCount
puts "count is: $count"

namespace eval engine {
  proc up {} {
    puts "<engine::up>"
  }
  proc down {} {
    puts "<engine::down>"
  }
}

engine::up
engine::down

array set arr {
  fst 1
  snd 2
  lst 3
}
parray arr
puts [array get arr]
puts [array names arr]

puts [info tclversion]
puts [info level]
puts [info cmdcount]
puts [info commands]
puts [info hostname]
puts [info nameofexecutable]
puts [info procs]
puts [info args greet]
puts [info body greet]

puts [::tcl::mathop::/ 7]
