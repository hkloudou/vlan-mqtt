include $(shell git rev-parse --show-toplevel)/include.mk

run:
	./main_linux_amd64 --vid=1 --mqtt=broker.emqx.io:1883 --cid=m000002 --ip=10.1.0.2/16
run3:
	./main_linux_amd64 --vid=1 --mqtt=broker.emqx.io:1883 --cid=m000003 --ip=10.1.0.3/16