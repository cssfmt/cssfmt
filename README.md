# CSSFMT

Inspired by gofmt and goimports.

# Usage

```
cssfmt: [flags] [path ...]
	-w write to source files, not to stdout
	-d display diffs instead of rewriting files
	-e show all errors, not the first 10 lines
	-v validate file silently and exit with code 1 if modification occurred
```

# Examples

```sh
$ cssfmt -d norm.css
--- norm.css.orig	2017-03-31 00:43:19.016151784 +0300
+++ norm.css	2017-03-31 00:43:19.016151784 +0300
@@ -211,7 +211,6 @@
 	-webkit-appearance: none;
 }
 
-
 ::-webkit-file-upload-button {
 	-webkit-appearance: button;
 	font: inherit;
@@ -237,3 +236,4 @@
 [hidden] {
 	display: none;
 }
+

$ cat norm.css | ./cssfmt -d
--- <standard input>.orig	2017-03-31 00:43:05.530800501 +0300
+++ <standard input>	2017-03-31 00:43:05.530800501 +0300
@@ -211,7 +211,6 @@
 	-webkit-appearance: none;
 }
 
-
 ::-webkit-file-upload-button {
 	-webkit-appearance: button;
 	font: inherit;
@@ -237,3 +236,4 @@
 [hidden] {
 	display: none;
 }
+

```
