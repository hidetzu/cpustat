all: bin/cpustat

.PHONY: bin/cpustat
bin/cpustat
	@docker build . --target bin --output bin/ --platform local
