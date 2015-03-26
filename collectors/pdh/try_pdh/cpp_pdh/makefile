all:
	make go
	make cxx

go:
	rm -f try_pdh.cpp
	go build -o pdh_go.exe

cxx:
	g++ -I.. -c wpdh.cpp -o wpdh.o
	cp try_pdh.cpp.bak try_pdh.cpp
	g++ -I.. try_pdh.cpp wpdh.o ../pdh.lib -lstdc++ -o pdh_cxx.exe	
	rm try_pdh.cpp

clean:
	rm -f *.exe wpdh.o try_pdh.cpp
