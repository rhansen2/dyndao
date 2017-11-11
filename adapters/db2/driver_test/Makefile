# dyndao makefile, just for testing for now

run:
	 go run main.go -conn='DATABASE=SAMPLE; HOSTNAME=$(DB2_HOST); PORT=$(DB2_PORT); PROTOCOL=TCPIP; UID=$(DB2_USER); PWD=$(DB2_PASS);'

test:
	DB2_DSN='DATABASE=SAMPLE; HOSTNAME=$(DB2_HOST); PORT=$(DB2_PORT); PROTOCOL=TCPIP; UID=$(DB2_USER); PWD=$(DB2_PASS);' go test -v 

cover:
	
	DB2_DSN='DATABASE=SAMPLE; HOSTNAME=$(DB2_HOST); PORT=$(DB2_PORT); PROTOCOL=TCPIP; UID=$(DB2_USER); PWD=$(DB2_PASS);' go test -cover
