ASWSG      = ../../src/aswsg
COMMON     = header.in footer.in
WEBPAGES   = a.html b.html c.html

all: $(WEBPAGES)

%.html : %.pg $(COMMON)
	$(ASWSG) $< > $@

c.html: part1.in part2.in part3.in

.PHONY: clean
clean:
	rm *.html
