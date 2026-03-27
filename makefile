all: bin/cpustat

.PHONY: bin/cpustat
bin/cpustat
	@docker build -f build/Dockerfile . --target bin --output bin/ --platform local
