include $(GOROOT)/src/Make.inc

TARG=dirdump

GOFILES=dirdump.go\
		pages.go\

include $(GOROOT)/src/Make.cmd
