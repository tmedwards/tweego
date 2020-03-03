<!-- ***********************************************************************************************
	FAQ & Tips
************************************************************************************************ -->
<h1 id="faq-and-tips">FAQ &amp; Tips</h1>

This is a collection of tips, from how to avoid pitfalls to best practices.

<p role="note"><b>Note:</b>
Suggestions for new entries may be submitted by <a href="https://github.com/tmedwards/tweego/issues">creating a new issue</a> at Tweego's <a href="https://github.com/tmedwards/tweego">source code repository</a>—though acceptance of submissions <strong><em>is not</em></strong> guaranteed.
</p>


<!-- ***************************************************************************
	Avoid processing files
**************************************************************************** -->
<span id="faq-and-tips-avoid-processing-files"></span>
## Avoid processing files

The way to avoid having Tweego process files is to not pass it the files in the first place—i.e., keep the files in question separate from the files you want Tweego to compile.

Using image files as an example, I would generally recommend a directory structure something like:

```
project_directory/
	images/
	src/
```

Where `src` is the directory you pass to Tweego, which only contains files you want it to compile—and possibly files that it will not process, like notes and whatnot.  For example, while within the project directory the command:

```
tweego -o project.html src
```

Will only compile the files in `src`, leaving the image files in `images` alone.


<!-- ***************************************************************************
	Convert Twee2 files to Twee v3
**************************************************************************** -->
<span id="faq-and-tips-convert-twee2-files-to-tweev3"></span>
## Convert Twee2 files to Twee&nbsp;v3

You may convert a Twee2 notation file to a Twee&nbsp;v3 notation file like so:

```
tweego -d -o twee_v3_file.twee twee2_file.tw2
```

Or, if the Twee2 notation file has a standard Twee file extension (`.tw`, `.twee`), like so:

```
tweego --twee2-compat -d -o twee_v3_file.twee twee2_file.twee
```
