package hickwall

type SystemInfo struct {
	Name                      string   // DST54236
	Domain                    string   // cn1.global.xxxx.com
	NumberOfProcessors        int      // 1
	NumberOfLogicalProcessors int      // 2
	Architecture              int      // 32, 64
	TotalPhsycialMemoryKb     int      // 12121121
	OS                        string   // Win7 Pro, Ubuntu, CentOS
	OSVersion                 string   // Service Pack 1 - 6.1.7601, 12.04, 6.5
	IPv4                      []string // [a, b]
}
