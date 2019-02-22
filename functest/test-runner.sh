#!/bin/sh
for testscript in $( ls ??-test-*.sh); do
	testname=$(basename -- "$testscript")
	testname="${testname%.*}"  # see http://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html

	if ./$testscript; then
		printf "%-64s: OK\n" "$testname"
	else
		if [ "$?" == "99" ] ; then
			printf "%-64s: SKIP\n" "$testname"
		else
			printf "%-64s: FAILED\n" "$testname"
			exit 1
		fi
	fi
done
exit 0
