Hint to embed c++ into go
===================

`extern "C"` in header file is the key. then remains are almost the same as C.

	#ifdef __cplusplus
	extern "C" {
	#endif
	
	int cpuusage();
	int getcpuload();
	
	#ifdef __cplusplus
	}
	#endif

