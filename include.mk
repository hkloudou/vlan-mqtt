.PHONY: default
.DEFAULT_GOAL := default
NAMESPACE = common
FP=common/private/vmqtt
CONFIGFILE=~/.kube/config_tc_fan
HELM_NAME=$(subst /,-,${FP}/${GIT_SUBPATH})
ALIYUN_BUCKET_NAME=zz-erp-files
ALIYUN_BUCKET_HOST=${ALIYUN_BUCKET_NAME}.oss-cn-hangzhou.aliyuncs.com
ALIYUN_BUCKET_DICTORY=/k8s-bin

# go install github.com/hkloudou/git-autotag@latest
ifneq ($(shell pwd),$(shell git rev-parse --show-toplevel))
	GIT_SUBPATH=$(subst $(shell git rev-parse --show-toplevel)/,,$(shell pwd))
	GIT_SUB_PARAME = -s ${GIT_SUBPATH}
	GIT_CLOSEDVERSION = $(shell git describe --abbrev=0  --match ${GIT_SUBPATH}/v[0-9]*\.[0-9]*\.[0-9]*)
else
	GIT_SUBPATH = main
	GIT_CLOSEDVERSION = $(shell git describe --abbrev=0  --match v[0-9]*\.[0-9]*\.[0-9]*)
endif
print:
	@echo sub: ${GIT_SUBPATH}
	@echo close: ${GIT_CLOSEDVERSION}
default:
	-git autotag -commit 'modify ${GIT_SUBPATH}' -f -p ${GIT_SUB_PARAME}
	@echo current version:`git describe`
git:
	- git autotag -commit 'auto commit ${GIT_SUBPATH}' -t -f -i -p ${GIT_SUB_PARAME}
	@echo current version:`git describe`
retag:
	-git autotag -commit 'retag $(GIT_CLOSEDVERSION)' -t -f -p ${GIT_SUB_PARAME}
	@echo current version:`git describe`
git-minor:
	git autotag -commit 'auto commit ${GIT_SUBPATH}' -t -f -i -p -l minor ${GIT_SUB_PARAME}
b:
	echo ${GIT_SUBPATH}
	mkdir -p $(shell git rev-parse --show-toplevel)/bin/${dir ${GIT_SUBPATH}}
	GOOS=linux GOARCH=amd64 sh ${shell go env GOMODCACHE}/github.com/hkloudou/xlib@v1.0.62/scripts/gobuild.sh $(shell basename $(GIT_SUBPATH)) $(shell git rev-parse --show-toplevel)/bin/${GIT_SUBPATH}_linux_amd64
	ossutilmac64 -c ~/.ossutilconfig cp -f $(shell git rev-parse --show-toplevel)/bin/${GIT_SUBPATH}_linux_amd64 oss://${ALIYUN_BUCKET_NAME}${ALIYUN_BUCKET_DICTORY}/${FP}/${GIT_SUBPATH}_linux_amd64  --snapshot-path=$(shell git rev-parse --show-toplevel)/__upload_log
	rm -rf $(shell git rev-parse --show-toplevel)/bin/${${GIT_SUBPATH}}*
u:
	helm del --kube-insecure-skip-tls-verify --kubeconfig $(CONFIGFILE)  ${HELM_NAME} -n $(NAMESPACE)
log:
	kubectl --kubeconfig $(CONFIGFILE)  logs -c basic -f ${shell kubectl --kubeconfig $(CONFIGFILE) get pods -n $(NAMESPACE) -o name | grep ${HELM_NAME}} -n $(NAMESPACE)
logi:
	kubectl --kubeconfig $(CONFIGFILE)  logs -c init-data -f ${shell kubectl --kubeconfig $(CONFIGFILE) get pods -n $(NAMESPACE) -o name| grep ${HELM_NAME}} -n $(NAMESPACE)
pwd:
	@echo ${shell kubectl --kubeconfig $(CONFIGFILE) get secret --namespace common mysql -o jsonpath="{.data.mysql-root-password}" | base64 -d}