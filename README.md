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

## Example

```bash
rebyre on  main via 🐹 v1.13.15 
➜ rebyre solve example_input.boole
Starting resolution:
Found an empty clause !!

Solution #0

(  )┬( x )┬( x | !c )┬( x | !d | !c )┬( x | !d | !a )
    |     |          |               └( !c | a | !d )
    |     |          └( x | d | !c )┬( x | d | a )
    |     |                         └( d | !c | !a )
    |     └( c | x )┬( c | y )┬( c | !z | y )
    |               |         └( z )
    |               └( x | !y )┬( x | !b | !y )
    |                          └( b )
    └( !x )┬( !c | !x )┬( !c | !y )┬( !c | !y | !z )
           |           |           └( z )
           |           └( y | !x )┬( y | !b | !x )
           |                      └( b )
           └( !x | c )┬( !x | !a | c )┬( !x | !a | d )
                      |               └( !a | !d | c )
                      └( a | c )┬( !d | a | c )
                                └( d | c | a )
```