# Data Serialization Toolbox (ds)

Semantic Diffs of Data Serialization languages.

## Current Capabilities

Parsing and diffs between/of:
- yaml
- json
- toml
- xml. Diffs with XML are inherently a little lossy, see the **XML Diffs** section.

Diffs can be printed in:
- go-syntax format
- A buggy syntax agnostic format. It can result in unexpected type coercion.

## Future Goals

Printing diffs in:
- A stable syntax agnostic format
- The relevant language

Other goals:
- Collapsing unchanged stanzas to make the output easier to understand

## Non-Goals

- ini support

## XML Diffs

XML diffs are a little lossy. The order of elements in an XML document is significant, but go maps are used
to store the elements. `ds diff` will treat xml like the following two as identical, despite them being
technically different:

```
<?xml version="1.0" encoding="UTF-8"?>
<note>
    <test2>value2</test2>
    <test>value</test>
</note>
```

```
<?xml version="1.0" encoding="UTF-8"?>
<note>
    <test>value</test>
    <test2>value2</test2>
</note>
```

If you *need* order-sensitive diffs of xml files, look into [xmldiff](https://pypi.org/project/xmldiff).

## Inspiration

This is partially a reimplementation, partially a generalization of:
- The phenomenal [json-diff](https://github.com/andreyvit/json-diff) by Andrey Tarantsov.
  - This one inspired a lot of the interface and features.
- The equally awesome [yamldiff](https://github.com/sahilm/yamldiff) by Sahil Muthoo.
  - This one showed me out some awesome libraries that helped write `dsdiff`.
