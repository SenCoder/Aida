.PHONY : clean all

all : Aida AidaService

clean :
	-rm -r Aida AidaService

#Aida : src/Aida/client/*.go
#	go build -o Aida Aida/client

Aida : src/Aida/client/*.cpp
	gcc src/Aida/client/socket_client.cpp -o Aida -lpthread

AidaService : src/Aida/service/*.go src/Aida/service/cloud/*.go
	go build -o AidaService Aida/service
