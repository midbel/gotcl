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

clock seconds

puts "count outside testUpvar: $count"
if { ::tcl::mathop::> 10 20 } then {
  puts "ok"
} elseif { ::tcl::mathop::!= 10 10 } then {
  puts "<equal>"
}

switch foobar {
  {foo*[a-z]*} { puts "foo" }
  bar { puts "bar" }
  default { puts "something else" }
}

puts "procs: [info procs]"
puts [info cmdcount]
puts [info level]

set now [clock seconds]
clock format $now "%Y-%m-%d %H:%M"
clock scan "2022-08-22 20:01" "%Y-%m-%d %H:%M"
