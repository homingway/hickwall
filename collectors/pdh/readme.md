this collector only works with `go1.4.2+`. 

known issues:
---------

* if build with `go1.3rc2`, once you `AddEnglishCounter` with **an invalid path, it will crash the application**, and no where to catch the exception!. `PdhValidatePath` also has the same issue. `go1.4.2` don't have this problem.

      // add an invalid counter path. this will crash application. 
      // with in cygwin env, it shows `segmentation fault`. 
      // under windows, appcrash event log will be generated in application log.
      AddEnglishCounter("\\Processes(_Total)\\Working Set")

      // this is the valid one
      AddEnglishCounter("\\Process(_Total)\\Working Set")


notes
-------
After encountered the `segmentation fault`, I tried a lot of ways to fix it:

* tried to use c to implement all functions, and wrap them with cgo. the bug remains. 
* tried to use c++ to catch exceptions, but always failed while compiling. 
* tried to use signal in c to catch SIGSEGV,SIGABRT, SIGILL, ... but bug remains. 
* accidentially upgraded `go1.3rc2` to `go1.4.2`, bug fixed!

`try_pdh` contains demos to embed c and c++ into go!