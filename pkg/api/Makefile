SWAGGER_OUT ?= swagger.yaml

SWAGGER ?= ../../test/tools/build/swagger

validate: ${SWAGGER_OUT}
	$(SWAGGER) validate ${SWAGGER_OUT}

serve: ${SWAGGER_OUT}
	$(SWAGGER) serve -F redoc -p=8080 swagger.yaml

.PHONY: ${SWAGGER_OUT}
${SWAGGER_OUT}:
	# generate doesn't remove file on error
	rm -f ${SWAGGER_OUT}
	$(SWAGGER) generate spec -x github.com/sigstore/rekor -x github.com/docker/docker/api/types/network -o ${SWAGGER_OUT} -i tags.yaml -w ./ -m
