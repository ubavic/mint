# Mint

Mint is a document (meta)language and preprocessor which tries to be simple and flexible.

Like in LaTex, in Mint grouping of text is done with braces. Unlike TeX, in Mint special character is not `\` but `@`, hence every command starts with `@` (eg. `@title`, `@bold`...). 

Mint hasn't any predefined command (scaping sequences `@@`, `@{` and `@}` only look like commands), not even basic document commands like that for paragraph, title or text decorations. All commands should be defined by user in schema file called `mint.yaml`.  

## Usage

You have to provide path to `.atex` file and `.yaml` schema:

```
mint -in "file.atex" -schema "schema.yaml" [-target TargetName]
```

See `./example`

## TODO

Mint is still early in development phase, and here is a list of features that might be developed at some point:

 + command IDs
 + command parameters
 + intuitive handling of multiple files
 + begin/end commands
 + implicit command arguments
 + parameter typing
 + parameter description
 + more optimized tokenizer/parser/writer
 + schema validation
 + JSON input/output
 + stable API
 + allow schema written in `atex` 
 + WASM filters
 + language server and extensions for editor

## License

The code is released under the [MIT license](LICENSE).
