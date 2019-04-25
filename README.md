# Data Serialization Diff (dsdiff) 

Semantic Diffs of Data Serialization languages.

## Future Goals

Diffs between:
- toml
- xml?
- ini?
- protobuf??? Probably not, but maybe

Other things:
- Pretty printing in the relevant language

## Current Capabilities

- yaml
- json

Printing the output of a syntax-agnostic diff

## Inspiration

This is partially a reimplementation, partially a 
- The phenomenal [json-diff](https://github.com/andreyvit/json-diff) by Andrey Tarantsov.
  - This one inspired a lot of the interface and features.
- The equally awesome [yamldiff](https://github.com/sahilm/yamldiff) by Sahil Muthoo.
  - This one showed me out some awesome libraries that handle most of the logic in `dsdiff`.
