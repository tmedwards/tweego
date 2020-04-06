<!-- ***********************************************************************************************
	Twee Notation
************************************************************************************************ -->
<h1 id="twee-notation">Twee Notation</h1>

In Twee and Twine, stories are arranged into units called passages.  Each passage has a name, optional attributes, and content.

There are two official Twee notations, Twee&nbsp;v3 and Twee&nbsp;v1, and an unofficial Twee2 notation.

* Twee&nbsp;v3 is the current official notation—see the <a href="https://github.com/iftechfoundation/twine-specs/blob/master/twee-3-specification.md" target="&#95;blank">twee-3-specification.md</a> for more information.
* Twee&nbsp;v1 is the classic/legacy official notation, which is a compatible subset of Twee&nbsp;v3.
* The unofficial Twee2 notation is primarily generated and used by the Twee2 compiler, which is largely abandonware.

By default, Tweego supports compiling from both of the official Twee notations and decompiling to Twee&nbsp;v3.  Compiling from the unofficial Twee2 notation is also supported via a compatibility mode, but is not enabled by default.  To load files with the Twee2 compatibility mode enabled, either the files must have a Twee2 extension (`.tw2`, `.twee2`) or its command line option (<kbd>--twee2-compat</kbd>) must be used.

<p class="warning" role="note"><b>Warning:</b>
It is <strong><em>strongly recommended</em></strong> that you do not enable Twee2 compatibility mode unless you absolutely need it.
</p>


<!-- ***************************************************************************
	Twee v3 Notation
**************************************************************************** -->
<span id="twee-notation-tweev3"></span>
## Twee&nbsp;v3 Notation

In the Twee&nbsp;v3 notation, passages consist of a passage declaration and a following content section.

A passage declaration must be a single line and is composed of the following components *(in order)*:

1. A required start token that must begin the line.  It is composed of a double colon (`::`).
2. A required passage name.
3. An optional tags block that must directly follow the passage name.  It is composed of a left square bracket (`[`), a space separated list of tags, and a right square bracket (`]`).
4. An optional metadata block that must directly follow either the tag block or, if the tag block is omitted, the passage name.  It is composed of an inline JSON chunk containing the optional properties `position` and `size`.

The passage content section begins with the very next line and continues until the next passage declaration.

<p class="tip" role="note"><b>Tip:</b>
For the sake of readability, it is recommended that each component within the passage declaration after the start token be preceded by one or more spaces and that, at least, one blank line is added between passages.
</p>

<p role="note"><b>Note:</b>
You will likely never need to create metadata blocks yourself.  When compiling, any missing metadata will be automatically generated for the compiled file.  When decompiling, they'll be automatically pulled from the compiled file.
</p>

<!-- *********************************************************************** -->

<span id="twee-notation-tweev3-passage-and-tag-name-escaping"></span>
### Passage And Tag Name Escaping

To prevent ambiguity during parsing, passage and tag names that include the optional tag or metadata block delimiters (`[`, `]`, `{`, `}`) must escape them.  The escapement mechanism is to prefix the escaped characters with a backslash (`\`).  Further, to avoid ambiguity with the escape character itself, non-escape backslashes must also be escaped via the same mechanism—e.g., `foo\bar` should be escaped as `foo\\bar`.

<p class="tip" role="note"><b>Tip:</b>
It is <strong><em>strongly recommended</em></strong> that you simply avoid needing to escape characters by not using the optional tag or metadata block delimiters within passage and tag names.
</p>

<p class="tip" role="note"><b>Tip:</b>
For different reasons, it is also <strong><em>strongly recommended</em></strong> that you avoid the use of the link markup separator delimiters (<code>|</code>, <code>-&gt;</code>, <code>&lt;-</code>) within passage and tag names.
</p>

<!-- *********************************************************************** -->

<span id="twee-notation-tweev3-example"></span>
### Example

#### Without any passage metadata

Exactly the same as Twee&nbsp;v1, save for the [Passage And Tag Name Escaping](#twee-notation-tweev3-passage-and-tag-name-escaping) rules.

```
:: A passage with no tags
Content of the "A passage with no tags" passage.


:: A tagged passage with three tags [alfa bravo charlie]
Content of the "A tagged passage with three tags" passage.
The three tags are: alfa, bravo, charlie.
```

#### With some passage metadata

Mostly likely to come from decompiling Twine&nbsp;2 or Twine&nbsp;1 compiled HTML files.

```
:: A passage with no tags {"position":"860,401"}
Content of the "A passage with no tags" passage.


:: A tagged passage with three tags [alfa bravo charlie] {"position":"860,530"}
Content of the "A tagged passage with three tags" passage.
The three tags are: alfa, bravo, charlie.
```


<!-- ***************************************************************************
	Twee v1 Notation
**************************************************************************** -->
<span id="twee-notation-tweev1"></span>
## Twee&nbsp;v1 Notation

<p class="warning" role="note"><b>Warning:</b>
Except in instances where you plan to interoperate with Twine&nbsp;1, it is <strong><em>strongly recommended</em></strong> that you do not create new files using the Twee&nbsp;v1 notation.  You should prefer the <a href="#twee-notation-tweev3">Twee&nbsp;v3 notation</a> instead.
</p>

Twee&nbsp;v1 notation is a subset of Twee&nbsp;v3 that lacks support for both the optional metadata block within passage declarations and passage and tag name escaping.

<!-- *********************************************************************** -->

<span id="twee-notation-tweev1-example"></span>
### Example

```
:: A passage with no tags
Content of the "A passage with no tags" passage.


:: A tagged passage with three tags [alfa bravo charlie]
Content of the "A tagged passage with three tags" passage.
The three tags are: alfa, bravo, charlie.
```


<!-- ***************************************************************************
	Twee2 Notation
**************************************************************************** -->
<span id="twee-notation-twee2"></span>
## Twee2 Notation

<p class="warning" role="note"><b>Warning:</b>
It is <strong><em>strongly recommended</em></strong> that you do not create new files using the unofficial Twee2 notation.  You should prefer the <a href="#twee-notation-tweev3">Twee&nbsp;v3 notation</a> instead.
</p>

The unofficial Twee2 notation is mostly identical to the Twee&nbsp;v1 notation, save that the passage declaration may also include an optional position block that must directly follow either the tag block or, if the tag block is omitted, the passage name.


<!-- *********************************************************************** -->

<span id="twee-notation-tweev1-example"></span>
### Example

```
:: A passage with no tags <860,401>
Content of the "A passage with no tags" passage.


:: A tagged passage with three tags [alfa bravo charlie] <860,530>
Content of the "A tagged passage with three tags" passage.
The three tags are: alfa, bravo, charlie.
```
