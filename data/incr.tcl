#! /usr/bin/tclsh
set  i 0
incr i # returns 1 and update the value of $i to 1
puts "result: $i (should be 1)" # print 1 to stdout
expr 10+10
echo "print from echo"
puts "argv: $argv"
puts "script [expr $i+$i*$tcl_command/$tcl_depth]"
if { ::tcl::mathop::< 10 20 } {
  puts "ok"
}
set count 0
for { set i 0 } { ::tcl::mathop::< $i 3 } { incr i } {
  puts "counter: [incr count] - value: $i"
}
puts "result: $i (should be 10)" # print 1 to stdout
puts "after for: $count"
set ne ::tcl::mathop::!=
while { $ne $i 0 } {
  puts "counter: [incr count] - value: [decr i]"
}
puts "after while: $count"
proc sayHello {who {message hello}} {
  puts "$message $who"
}
sayHello foobar
sayHello foobar goodbye
clock seconds

proc testUpvar {} {
  upvar 1 count counter
  incr counter
  puts "results counter (testUpvar): $counter"
}

testUpvar
puts "count outside testUpvar: $count"
if { ::tcl::mathop::> 10 20 } then {
  puts "ok"
} elseif { ::tcl::mathop::!= 10 10 } then {
  puts "<equal>"
}
