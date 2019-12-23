<!-- ***********************************************************************************************
	Getting Started
************************************************************************************************ -->
<h1 id="getting-started">Getting Started</h1>


<!-- ***************************************************************************
	Overview
**************************************************************************** -->
<span id="getting-started-overview"></span>
## Overview

<p class="tip" role="note"><b>Tip:</b>
In practice, most settings will be handled either by story configuration or via the command line, so the only configuration step that's absolutely necessary to begin using Tweego is to enable it to find your story formats.
</p>

Tweego may be configured in a variety of ways—by environment variable, story configuration, and command line options.

The various methods for specifying configuration settings cascade in the following order:

1. Program defaults.
2. Environment variables.
3. Story configuration.  See the [`StoryData` passage](#special-passages-storydata) for more information.
4. Command line.  See [Usage](#usage) for more information.


<!-- ***************************************************************************
	Program Defaults
**************************************************************************** -->
<span id="getting-started-program-defaults"></span>
## Program Defaults

<dl>
<dt>Input charset</dt>
<dd>
	<p>The default character set is <code>utf-8</code>, failing over to <code>windows-1252</code> if the input files are not in UTF-8.</p>
	<p class="tip" role="note"><b>Tip:</b> It is <strong><em>strongly recommended</em></strong> that you use UTF-8 for all of your text files.</p>
</dd>
<dt>Story format</dt>
	<dd>The default story format (by ID) is <code>sugarcube-2</code>.</dd>
<dt>Output file</dt>
	<dd>The default output file is <code>-</code>, which is shorthand for <a href="https://en.wikipedia.org/wiki/Standard_streams" target="&#95;blank"><i>standard output</i></a>.</dd>
<dt>Starting passage</dt>
	<dd>The default starting passage name is <code>Start</code>.</dd>
</dl>


<!-- ***************************************************************************
	Environment Variables
**************************************************************************** -->
<span id="getting-started-environment-variables"></span>
## Environment Variables

<dl>
<dt id="getting-started-environment-variables-tweego-path"><var>TWEEGO_PATH</var></dt>
<dd>
	<p>Path(s) to search for story formats.  The value should be a list of directories to search for story formats.  You may specify one directory or several.  The format is exactly the same as any other <em>path type</em> environment variable for your operating system.</p>
	<p class="tip" role="note"><b>Tip:</b> Setting <var>TWEEGO_PATH</var> is only necessary if you intend to place your story formats outside of the directories normally searched by Tweego.  See <a href="#getting-started-story-formats-search-directories">Search Directories</a> for more information.</p>
	<p role="note"><b>Note:</b> To separate multiple directories within <em>path</em> variables, Unix/Unix-like operating systems use the colon (<kbd>:</kbd>), while Windows uses the semi-colon (<kbd>;</kbd>).  Only relevant if you intend to specify multiple directories.</p>
	<p><strong>Unix/Unix-like examples</strong></p>
	<p>If you wanted Tweego to search <code>/usr/local/storyformats</code>, then you'd set <code>TWEEGO_PATH</code> to:</p>
	<pre><code>/usr/local/storyformats</code></pre>
	<p>If you wanted Tweego to search <code>/storyformats</code> and <code>/usr/local/storyformats</code>, then you'd set <code>TWEEGO_PATH</code> to:</p>
	<pre><code>/storyformats:/usr/local/storyformats</code></pre>
	<p><strong>Windows examples</strong></p>
	<p>If you wanted Tweego to search <code>C:\\storyformats</code>, then you'd set <code>TWEEGO_PATH</code> to:</p>
	<pre><code>C:\storyformats</code></pre>
	<p>If you wanted Tweego to search <code>C:\storyformats</code> and <code>D:\storyformats</code>, then you'd set <code>TWEEGO_PATH</code> to:</p>
	<pre><code>C:\storyformats;D:\storyformats</code></pre>
</dd>
</dl>


<!-- ***************************************************************************
	Story Formats
**************************************************************************** -->
<span id="getting-started-story-formats"></span>
## Story Formats

<p role="note"><b>Note:</b>
Throughout this document the terms <code>story format</code> and <code>format</code> are virtually always used to encompass both story and proofing formats.
</p>

Tweego should be compatible with *all* story formats—i.e., those written for Twine&nbsp;2, Twine&nbsp;1 ≥v1.4.0, and Twine&nbsp;1 ≤v1.3.5.

Installing a story format can be as simple as moving its directory into one of the directories Tweego searches for story formats—see [Search Directories](#getting-started-story-formats-search-directories) for more information.  Each installed story format, which includes separate versions of the same story format, should have its own <em>unique</em> directory within your story formats directory—i.e., if you have both SugarCube v2 and v1 installed, then they should each have their own separate directory; e.g., `sugarcube-2` and `sugarcube-1`.  Do not create additional sub-directories, combine directories, or rename a story format's files.

<p class="tip" role="note"><b>Tip:</b>
To ensure a story format has been installed correctly, use the list-formats command line option (<kbd>--list-formats</kbd>) to see if Tweego lists it as an available format.
</p>

<p class="warning" role="note"><b>Warning:</b>
Twine&nbsp;2 story formats are, ostensibly, encoded as JSON-P.  Unfortunately, some story formats deviate from proper JSON encoding and are thus broken.  Tweego uses a strict JSON decoder and cannot decode such broken story formats for use.  Should you receive a story format decoding error, all reports should go to the format's developer.
</p>

<!-- *********************************************************************** -->

<span id="getting-started-story-formats-search-directories"></span>
### Search Directories

When Tweego is run, it finds story formats to use by searching the following directories: *(in order)*

1. The directories <kbd>storyformats</kbd> and <kbd>.storyformats</kbd> within its <em>program directory</em>—i.e., the directory where Tweego's binary file is located.
2. The directories <kbd>storyformats</kbd> and <kbd>.storyformats</kbd> within the <em>user's home directory</em>—i.e., either the value of the <var>HOME</var> environment variable or the operating system specified home directory.
3. The directories <kbd>storyformats</kbd> and <kbd>.storyformats</kbd> within the <em>current working directory</em>—i.e., the directory that you are executing Tweego from.
4. The directories specified via the <var>TWEEGO_PATH</var> environment variable.  See <a href="#getting-started-environment-variables">Environment Variables</a> for more information.

<p role="note"><b>Note:</b>
For legacy compatibility, the following directories are also checked during steps #1–3: <kbd>story-formats</kbd>, <kbd>storyFormats</kbd>, and <kbd>targets</kbd>.  You are <strong><em>strongly encouraged</em></strong> to use one of the directory names listed above instead.
</p>

<p class="warning" role="note"><b>Warning:</b>
A story format's directory name is used as its <strong><em>unique</em></strong> ID within the story format list.  As a consequence, if multiple story formats, from different search paths, have the same directory name, then only the last one found will be registered.
</p>
