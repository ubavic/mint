# Mint

Mint is a document (meta)language and preprocessor designed to be simple and flexible.

Like in LaTeX, grouping of text in Mint is done using braces. However, unlike TeX, the special character in Mint is not `\` but `@`. Therefore, every command starts with `@` (e.g., `@title`, `@bold`...).

Mint doesn’t have any predefined commands (the escape sequences `@@`, `@{`, and `@}` only resemble commands). Even basic document commands like those for paragraphs, titles, or text decorations are not predefined. All commands must be defined by the user in a YAML schema file.

## Usage

You have to provide path to `.atex` file and `.yaml` schema:

```
mint -in "file.atex" -schema "schema.yaml" [-target TargetName]
```

See `./example`

## TODO

Mint is still in the early development phase. Below is a list of features that may be developed in the future:

 + Command IDs
 + Command parameters
 + Intuitive handling of multiple files
 + Begin/end commands
 + Implicit command arguments
 + Parameter typing
 + Parameter description
 + More optimized tokenizer/parser/writer
 + Schema validation
 + JSON input/output
 + Stable API
 + Allow schema written in atex
 + WASM filters
 + Language server and extensions for editors

## License

The code is released under the [MIT license](LICENSE).
