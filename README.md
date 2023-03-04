# shelldiff

```go
func Diff(s1 string, s2 string, sw io.StringWriter) error
```

Prints a diff of two shell script snippets.

The diff uses comments in the scripts to identify matching sections.

For example:

comparing:
```bash
# get config
C=...

# compute
R=C+D

# print result
echo $R
```

with:
```bash
# compute
R=C-D
```

will print the following diff:

```patch
-[get config] C=...
-[compute] R=C+D
+[compute] R=C-D
-[print result] echo $R
```

This is likely to be more meaningful than a general line by line diff which has no understanding of the semantic
structure of the script.
