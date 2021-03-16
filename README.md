# A simple file checksum and comparison tool written in GO
This simple tool will calculate and compare the file checksum (using MD5) between two folders.
The checksums will be also saved to file for references.

## Usage
checksum-tool -input1 \<path> -input2 \<path>

## Console output
Running...\
==Mismatched checksum==\
\<file list>\
==Input 1 file(s) not found in input 2==\
\<file list>\
==Input 2 file(s) not found in input 1==\
\<file list>\
Done. Duration: \<xx>ms