# Introduce
iostat implement in Golang.

# How to use?
1. go get github.com/DoLolitingyu/goiostat
2. iostat.GetData()

# How's the effect? With iostat cmd?
You can run iostat_test.go to verify this. 

99% close you use iostat cmd.

1. 
goiostat: {vda 0 1 0 3 0 15.99 10.67 0 0 0 0 0 0}
iostat:   [vda 0.00 1.00 0.00 3.00 0.00 16.00 10.67 0.00 0.00 0.00 0.00 0.00 0.00 ]

2.
{vda 0 1 0 7 0 31.98 9.14 0 0 0 0 0 0}
[vda 0.00 1.00 0.00 7.00 0.00 32.00 9.14 0.00 0.00 0.00 0.00 0.00 0.00 ]
