#!/usr/bin/bash

mkdir -p foo/bar/baz/
mkdir -p foo/bar/boo/

touch foo/tmp1.txt
touch foo/tmp2.txt

touch foo/bar/temp1.txt
rm -f foo/bar/buzz
ln -s ${PWD}/foo/tmp1.txt foo/bar/buzz

touch foo/bar/baz/t1.txt
touch foo/bar/baz/t2.txt
touch foo/bar/baz/t3.txt
rm -f foo/bar/baz/broken_sl
ln -s noexist foo/bar/baz/broken_sl

touch foo/bar/boo/rt.log
touch foo/bar/boo/tr.log
