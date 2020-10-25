# rebyre

Just a small tool I wrote for doing refutiation by resolution, on a proposition in CNF.

## Usage

Call the `solve` command and point it to an input file, the rest is magic (not really).

```bash
$ rebyre solve example_input.boole
```

Input has to be written in CNF, logical `and` as a `&` and logical `or` as a `|`.
For more info regarding the input format, view the example_input.boole file.

If you want to see more details on the resolution process the program does, add the `verbose` flag. It then prints out each clause that it finds, together with an id and the clauses that were used to derive the clause.

```
id | name            | id of clause a | id of clause b
1  | ( x | !d | !a ) |                |
2  | ( !c | a | !d ) |                |
3  | ( c | !z | y )  |                |
...| ...             | ...            |
16 | ( x | !d | !c ) | 1              | 2


```