# MINT

MINT is a document markup language designed to be extensible and informative for humans.

Like in LaTeX, grouping of text in Mint is done using braces. However, unlike TeX, the special character in Mint is not `\` but `@`. Therefore, every command starts with `@` (e.g., `@title`, `@bold`...).

Mint doesn’t have any predefined commands (the escape sequences `@@`, `@{`, and `@}` only resemble commands). Even basic document commands like those for paragraphs, titles, or text decorations are not predefined. All commands must be defined by the user in a YAML schema file.

## Usage

You have to provide path to `.mint` file and `.yaml` schema:

```
mint -in INPUT_FILE -schema SCHEMA_FILE [-target TARGET] [-out OUTPUT_FILE]
```

## TODO

Mint is still in the early development phase. Below is a list of features that may be developed in the future:

- Parameter typing
- Command parameters
- Begin/end commands
- More optimized tokenizer/parser/writer
- Stable API
- Language server and extensions for editors
- Schema written in MINT

## License

The code is released under the [MIT license](LICENSE).
