# Introduce
sysstat implement in Golang.
include iostat and mpstat.

Don't have sysstat tool in many arm linux, so do this project, you can use it immediately.

# How to use?
1. go get github.com/DoLolitingyu/goiostat
2. iostat.GetData()
3. mpstat.GetData()

# How's the effect? With iostat cmd?
You can run iostat_test.go to verify this. 

99.9999999% close you use iostat cmd.
Very small errors are also due to timing statistics, and the commands themselves are not exact.

1. 
- goiostat:{vda 0 1 0 3 0 15.99 10.67 0 0 0 0 0 0}
- iostat: [vda 0.00 1.00 0.00 3.00 0.00 16.00 10.67 0.00 0.00 0.00 0.00 0.00 0.00 ]

2.
- gompstat: CPU,%usr,%nice,%sys,%iowait,%irq,%soft,%steal,%guest,%gnice,%idle
go-cpu,36.21,0.00,10.34,0.57,0.00,0.00,2.30,0.00,0.00,50.57
go-cpu0,32.58,0.00,13.48,2.25,0.00,0.00,2.25,0.00,0.00,49.44
go-cpu1,39.53,0.00,6.98,0.00,0.00,0.00,2.33,0.00,0.00,51.16

- mpstat: CPU,%usr,%nice,%sys,%iowait,%irq,%soft,%steal,%guest,%gnice,%idle
cpu-all,36.36,0.00,10.23,0.57,0.00,0.00,2.27,0.00,0.00,50.57
cpu-0,32.22,0.00,13.33,2.22,0.00,0.00,2.22,0.00,0.00,50.00
cpu-1,39.53,0.00,6.98,0.00,0.00,0.00,2.33,0.00,0.00,51.16
