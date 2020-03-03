<!-- ***********************************************************************************************
	Usage
************************************************************************************************ -->
<h1 id="usage">Usage</h1>


<!-- ***************************************************************************
	Overview
**************************************************************************** -->
<span id="usage-overview"></span>
## Overview

<p class="tip" role="note"><b>Tip:</b>
At any time you may pass the help option (<kbd>-h</kbd>, <kbd>--help</kbd>) to Tweego to show its built-in help.
</p>

Basic command line usage is as follows:

```
tweego [options] sources…
```

Where <code>[options]</code> are mostly optional configuration flags—see [Options](#usage-options)—and <code>sources</code> are the input sources which may consist of supported files and/or directories to recursively search for such files.  Many types of files are supported as input sources—see [Supported Files](#usage-supported-files) for more information.


<!-- ***************************************************************************
	Options
**************************************************************************** -->
<span id="usage-options"></span>
## Options

<dl>
<dt><kbd>-a</kbd>, <kbd>--archive-twine2</kbd></dt><dd>Output Twine&nbsp;2 archive, instead of compiled HTML.</dd>
<dt><kbd>--archive-twine1</kbd></dt><dd>Output Twine&nbsp;1 archive, instead of compiled HTML.</dd>
<dt><kbd>-c SET</kbd>, <kbd>--charset=SET</kbd></dt>
<dd>
	<p>Name of the input character set (default: <code>"utf-8"</code>, fallback: <code>"windows-1252"</code>).  Necessary only if the input files are not in either UTF-8 or the fallback character set.</p>
	<p class="tip" role="note"><b>Tip:</b> It is <strong><em>strongly recommended</em></strong> that you use UTF-8 for all of your text files.</p>
</dd>
<dt><kbd>-d</kbd>, <kbd>--decompile-twee3</kbd></dt><dd>Output Twee 3 source code, instead of compiled HTML.  See <a href="#twee-notation-tweev3">Twee&nbsp;v3 Notation</a> for more information.</dd>
<dt><kbd>--decompile-twee1</kbd></dt>
<dd>
	<p>Output Twee 1 source code, instead of compiled HTML.  See <a href="#twee-notation-tweev1">Twee&nbsp;v1 Notation</a> for more information.</p>
	<p role="note"><b>Note:</b> Except in instances where you plan to interoperate with Twine&nbsp;1, it is <strong><em>strongly recommended</em></strong> that you decompile to Twee&nbsp;v3 notation rather than Twee&nbsp;v1.</p>
</dd>
<dt><kbd>-f NAME</kbd>, <kbd>--format=NAME</kbd></dt><dd>ID of the story format (default: <code>"sugarcube-2"</code>).</dd>
<dt><kbd>-h</kbd>, <kbd>--help</kbd></dt><dd>Print the built-in help, then exit.</dd>
<dt><kbd>--head=FILE</kbd></dt><dd>Name of the file whose contents will be appended as-is to the &lt;head&gt; element of the compiled HTML.</dd>
<dt><kbd>--list-charsets</kbd></dt><dd>List the supported input character sets, then exit.</dd>
<dt><kbd>--list-formats</kbd></dt><dd>List the available story formats, then exit.</dd>
<dt><kbd>--log-files</kbd></dt>
<dd>
	<p>Log the processed input files.</p>
	<p role="note"><b>Note:</b> Unsupported when watch mode (<kbd>-w</kbd>, <kbd>--watch</kbd>) is enabled.</p>
</dd>
<dt><kbd>-l</kbd>, <kbd>--log-stats</kbd></dt>
<dd>
	<p>Log various story statistics.  Primarily, passage and word counts.</p>
	<p role="note"><b>Note:</b> Unsupported when watch mode (<kbd>-w</kbd>, <kbd>--watch</kbd>) is enabled.</p>
</dd>
<dt><kbd>-m SRC</kbd>, <kbd>--module=SRC</kbd></dt><dd>Module sources (repeatable); may consist of supported files and/or directories to recursively search for such files.  Each file will be wrapped within the appropriate markup and bundled into the &lt;head&gt; element of the compiled HTML.  Supported files: <code>.css</code>, <code>.js</code>, <code>.otf</code>, <code>.ttf</code>, <code>.woff</code>, <code>.woff2</code>.</dd>
<dt><kbd>--no-trim</kbd></dt><dd>
	<p>Do not trim whitespace surrounding passages—i.e., whitespace preceding and trailing the actual text of the passage.  By default, such whitespace is removed when processing passages.</p>
	<p role="note"><b>Note:</b> It is recommended that you do not disable passage trimming.</p>
</dd>
<dt><kbd>-o FILE</kbd>, <kbd>--output=FILE</kbd></dt><dd>Name of the output file (default: <kbd>-</kbd>; i.e., <a href="https://en.wikipedia.org/wiki/Standard_streams" target="&#95;blank"><i>standard output</i></a>).</dd>
<dt><kbd>-s NAME</kbd>, <kbd>--start=NAME</kbd></dt><dd>Name of the starting passage (default: the passage set by the story data, elsewise <code>"Start"</code>).</dd>
<dt><kbd>-t</kbd>, <kbd>--test</kbd></dt><dd>Compile in test mode; only for story formats in the Twine&nbsp;2 style.</dd>
<dt><kbd>--twee2-compat</kbd></dt><dd>Enable Twee2 source compatibility mode; files with the <code>.tw2</code> or <code>.twee2</code> extensions automatically have compatibility mode enabled.</dd>
<dt><kbd>-v</kbd>, <kbd>--version</kbd></dt><dd>Print version information, then exit.</dd>
<dt><kbd>-w</kbd>, <kbd>--watch</kbd></dt><dd>Start watch mode; watch input sources for changes, rebuilding the output as necessary.</dd>
</dl>


<!-- ***************************************************************************
	Supported Files
**************************************************************************** -->
<span id="usage-supported-files"></span>
## Supported Files

Tweego supports various types of files for use in projects.  File types are recognized by filename extension, so all files ***must*** have an extension.

The following extensions are supported:

<dl>
<dt><code>.tw</code>, <code>.twee</code></dt>
<dd>
	<p>Twee notation source files to process for passages.</p>
	<p role="note"><b>Note:</b> If any of these files are in the unofficial Twee2 notation, you must manually enable the Twee2 compatibility mode via its command line option (<kbd>--twee2-compat</kbd>).</p>
</dd>
<dt><code>.tw2</code>, <code>.twee2</code></dt>
<dd>Unofficial Twee2 notation source files to process for passages.  Twee2 compatibility mode is automatically enabled for files with these extensions.</dd>
<dt><code>.htm</code>, <code>.html</code></dt>
<dd>HTML source files to process for passages, either compiled files or story archives.</dd>
<dt><code>.css</code></dt>
<dd>CSS source files to bundle.</dd>
<dt><code>.js</code></dt>
<dd>JavaScript source files to bundle.</dd>
<dt><code>.otf</code>, <code>.ttf</code>, <code>.woff</code>, <code>.woff2</code></dt>
<dd>Font files to bundle, as <code>@font-face</code> style rules.  The generated name of the font family will be the font's base filename sans its extension—e.g., the family name for <code>chinacat.tff</code> will be <code>chinacat</code>.</dd>
<dt><code>.gif</code>, <code>.jpeg</code>, <code>.jpg</code>, <code>.png</code>, <code>.svg</code>, <code>.tif</code>, <code>.tiff</code>, <code>.webp</code></dt>
<dd>
	<p>Image files to bundle, as image passages.  The generated name of the image passage will be the base filename sans its extension—e.g., the passage name for <code>rainboom.jpg</code> will be <code>rainboom</code>.</p>
	<p role="note"><b>Note:</b>
	As of this writing, image passages are only natively supported by SugarCube (all versions) and the Twine&nbsp;1 ≥v1.4 vanilla story formats.
	</p>
</dd>
<dt><code>.aac</code>, <code>.flac</code>, <code>.m4a</code>, <code>.mp3</code>, <code>.oga</code>, <code>.ogg</code>, <code>.opus</code>, <code>.wav</code>, <code>.wave</code>, <code>.weba</code></dt>
<dd>
	<p>Audio files to bundle, as audio passages.  The generated name of the audio passage will be the base filename sans its extension—e.g., the passage name for <code>swamped.mp3</code> will be <code>swamped</code>.</p>
	<p role="note"><b>Note:</b>
	As of this writing, audio passages are only natively supported by SugarCube ≥v2.24.0.
	</p>
</dd>
<dt><code>.mp4</code>, <code>.ogv</code>, <code>.webm</code></dt>
<dd>
	<p>Video files to bundle, as video passages.  The generated name of the video passage will be the base filename sans its extension—e.g., the passage name for <code>cutscene.mp4</code> will be <code>cutscene</code>.</p>
	<p role="note"><b>Note:</b>
	As of this writing, video passages are only natively supported by SugarCube ≥v2.24.0.
	</p>
</dd>
<dt><code>.vtt</code></dt>
<dd>
	<p>Text track files to bundle, as text track passages.  The generated name of the text track passage will be the base filename sans its extension—e.g., the passage name for <code>captions.vtt</code> will be <code>captions</code>.</p>
	<p role="note"><b>Note:</b>
	As of this writing, text track passages are only natively supported by SugarCube ≥v2.24.0.
	</p>
</dd>
</dl>


<!-- ***************************************************************************
	File & Directory Handling
**************************************************************************** -->
<span id="usage-file-and-directory-handling"></span>
## File &amp; Directory Handling

Tweego allows you to specify an arbitrary number of files and directories on the command line for processing.  In addition to those manually specified, it will recursively search all directories encountered looking for additional files and directories to process.  Generally, this means that you only have to specify the base source directory of your project and Tweego will find all of its files automatically.


<!-- ***************************************************************************
	Basic Examples
**************************************************************************** -->
<span id="usage-basic-examples"></span>
## Basic Examples

Compile <kbd>example_1.twee</kbd> as <kbd>example_1.html</kbd> with the default story format:

```
tweego -o example_1.html example_1.twee
```

Compile all files in <kbd>example_directory_2</kbd> as <kbd>example_2.html</kbd> with the default story format:

```
tweego -o example_2.html example_directory_2
```

Compile <kbd>example_3.twee</kbd> as <kbd>example_3.html</kbd> with the story format <kbd>snowman</kbd>:

```
tweego -f snowman -o example_3.html example_3.twee
```

Compile all files in <kbd>example_directory_4</kbd> as <kbd>example_4.html</kbd> with the default story format while also bundling all files in <kbd>modules_directory_4</kbd> into the &lt;head&gt; element of the compiled HTML:

```
tweego -o example_4.html -m modules_directory_4 example_directory_4
```

Decompile <kbd>example_5.html</kbd> as <kbd>example_5.twee</kbd>:

```
tweego -d -o example5.twee example5.html
```
