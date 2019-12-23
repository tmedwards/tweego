<!-- ***********************************************************************************************
	Special Passages & Tags
************************************************************************************************ -->
<h1 id="special">Special Passages &amp; Tags</h1>

Passages and tags that have special meaning to Tweego.

<p role="note"><b>Note:</b>
This is <em>not</em> a exhaustive list of all special passages and tags that may have meaning to story formats—or other compilers.  See the documentation of the specific story format—or compiler—for their list of special passages and tags.
</p>

<p class="warning" role="note"><b>Warning:</b>
The names of all special passages and tags listed herein are case sensitive, thus must be spelled <em>exactly</em> as shown.
</p>


<!-- ***************************************************************************
	Special Passages
**************************************************************************** -->
<span id="special-passages"></span>
## Special Passages

<!-- *********************************************************************** -->

<span id="special-passages-start"></span>
### `Start`

The `Start` passage will, by default, be used as the starting passage—i.e. the first normal passage displayed to the player.  That behavior may be overridden via either the <var>start</var> property from the [`StoryData` passage](#special-passages-storydata) or the start command line option (<kbd>-s NAME</kbd>, <kbd>--start=NAME</kbd>).

<p class="tip" role="note"><b>Tip:</b>
It is <strong><em>strongly recommended</em></strong> that you simply use the default starting name, <code>Start</code>, when beginning new projects.
</p>

<!-- *********************************************************************** -->

<span id="special-passages-storydata"></span>
### `StoryData`

The `StoryData` passage may be used to specify basic project settings.  Its contents must consist of a JSON chunk, which is, generally, pretty-printed—i.e., line-broken and indented.

The core properties used with all story formats include:

- <var>ifid</var>: (string) Required.  The project's Interactive Fiction IDentifier (IFID), which is a unique code used to identify your project—similar to the ISBN code assigned to a book.  If the project does not already have an IFID, you may omit the property and Tweego will automatically generate one for you with instructions on how to copy it into the chunk.
- <var>start</var>: (string) Optional.  The name of the starting passage.  If omitted, defaults to the [`Start` passage](#special-passages-start).

The properties used only with Twine&nbsp;2-style story formats include:

- <var>format</var>: (string) Optional.  The name of the story format to compile against—e.g., `SugarCube`, `Harlowe`, `Chapbook`, `Snowman`.
- <var>format-version</var>: (string) Optional.  The version of the story format to compile against—e.g., `2.29.0`.  From the installed story formats matching the name specified in <var>format</var>, Tweego will attempt to use the greatest version that matches the specified major version—i.e., if <var>format-version</var> is `2.0.0` and you have the versions `1.0.0`, `2.0.0`, `2.5.0`, and `3.0.0` installed, Tweego will choose `2.5.0`.

<p role="note"><b>Note:</b>
The above is <em>not</em> an exhaustive list of all Twine&nbsp;2-style story format properties.  There are others available that are only useful when actually interoperating with Twine&nbsp;2—e.g, <var>tag-colors</var> and <var>zoom</var>.  See the <a href="https://github.com/iftechfoundation/twine-specs/blob/master/twee-3-specification.md" target="&#95;blank">twee-3-specification.md</a> for more information.
</p>

<p class="tip" role="note"><b>Tip:</b>
To compile against a specific version of a story format, use the format command line option (<kbd>-f NAME</kbd>, <kbd>--format=NAME</kbd>) to specify the version by its ID.  If you don't know the ID, use the list-formats command line option (<kbd>--list-formats</kbd>) to find it.
</p>

<p class="warning" role="note"><b>Warning:</b>
JSON chunks are not JavaScript object literals, though they look much alike.  Property names must always be double quoted and you should not include a trailing comma after the last property.
</p>

#### Example

```
:: StoryData
{
	"ifid": "D674C58C-DEFA-4F70-B7A2-27742230C0FC",
	"format": "SugarCube",
	"format-version": "2.29.0",
	"start": "My Starting Passage"
}
```

<!-- *********************************************************************** -->

<span id="special-passages-storytitle"></span>
### `StoryTitle`

The contents of the `StoryTitle` passage will be used as the name/title of the story.


<!-- ***************************************************************************
	Special Tags
**************************************************************************** -->
<span id="special-tags"></span>
## Special Tags

<!-- *********************************************************************** -->

<span id="special-tags-script"></span>
### `script`

The `script` tag denotes that the passage's contents are JavaScript code.

<p role="note"><b>Note:</b>
In general, Tweego makes creating script passages unnecessary as it will automatically bundle any JavaScript source files (<code>.js</code>) it encounters into your project.
</p>

<!-- *********************************************************************** -->

<span id="special-tags-stylesheet"></span>
### `stylesheet`

The `stylesheet` tag denotes that the passage's contents are CSS rules.

<p role="note"><b>Note:</b>
In general, Tweego makes creating stylesheet passages unnecessary as it will automatically bundle any CSS source files (<code>.css</code>) it encounters into your project.
</p>
