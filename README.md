# Tweego

[Tweego](http://www.motoslave.net/tweego/) is a free (gratis and libre) command line compiler for [Twine/Twee](http://twinery.org/) story formats, written in [Go](http://golang.org/).

See [Tweego's documentation](http://www.motoslave.net/tweego/docs/) for information on how to set it up and use it.

Tweego has a Trello board, [Tweego TODO](https://trello.com/b/l5xuRzFu).  Feel free to comment.

## INSTALLATION

You may either download one of the precompiled binaries from [Tweego's website](http://www.motoslave.net/tweego/), which are available in both 32- and 64-bit versions for multiple operating systems, or build Tweego from source (see **BUILDING FROM SOURCE** below).

## BUILDING FROM SOURCE

If you want to build Tweego from scratch, rather than grabbing one of the precompiled binaries off of its website, then these instructions are for you.

Tweego is written in the Go programming language, so you'll need to install it, if you don't already have it.  Additionally, to retrieve Go packages—like Tweego and its dependencies—from source control repositories, you'll need to install Git.

1. [Download and install the Go programming language (want ≥v1.13)](http://golang.org/doc/install)
2. [Download and install the Git source control management tool](https://git-scm.com/)

Once all the tooling is installed and set up, the next step is to fetch the Tweego source code.  Open a shell to wherever you wish to store the code and run the following command:

```
git clone https://github.com/tmedwards/tweego.git
```

Next, change to the directory that the previous command created:

```
cd tweego
```

Next, fetch Tweego's dependencies via the following command:

```
go get
```

You should now have Tweego and all its dependencies downloaded, so you may now compile and install it to your `GOPATH` by running the following command:

```
go install
```

Assuming that completed with no errors, your compiled Tweego binary should be in your `GOPATH`'s `bin` directory.  To run Tweego you'll need to either have added your `GOPATH` `bin` to your `PATH` environment variable—this was likely done for you automatically during the installation of Go—or copy the binary to an existing directory within your `PATH`.

Alternatively.  If you just want to compile Tweego, so that you can manually copy the binary to wherever you wish, use the following command instead:

```
go build
```

Assuming that completed with no errors, your compiled Tweego binary should be in the current directory—likely named either `tweego` or `tweego.exe` depending on your OS.

Finally, see [Tweego's documentation](http://www.motoslave.net/tweego/docs/) for information on how to set it up and use it.
