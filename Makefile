# compile and load bin to remote host
all:
	export CGO_ENABLED=0 GOOS=linux GOARCH=amd64; go build; export GOOS="" GOARCH="" CGO_ENABLED=""
	gzip -f dnslite
	scp dnslite.gz hz:/root/
	rm -f dnslite.gz
