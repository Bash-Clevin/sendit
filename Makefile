build:
	@go build -o bin/server

run: build
	@./bin/server

ssh:
	@ssh-keygen -f "" -R [localhost]:2222
	@cat big.txt | ssh localhost -p 2222

scp:
	ssh-keygen -f "" -R [localhost]:2222
	@scp -P 2222 main.go localhost:aa@dd.com

scpfolder:
	@ssh-keygen -f "" -R [localhost]:2222
	@scp -P 2222 -r testfolder localhost:aa@dd.com

badscp:
	@ssh-keygen -f "" -R [localhost]:2222
	@scp -P 2222 main.go localhost:aaa.com