# go-xtp - An experimental Go and MoonBit code generator for XTP Extension Plugins

For more information about XTP, please see:
https://docs.xtp.dylibso.com/docs/overview

## xtp2code

To install the `xtp2code` binary, first make sure the [Go] programming language
is installed, and then run:

```bash
$ go install github.com/gmlewis/go-xtp/cmd/xtp2code@latest
```

To check what version of `xtp2code` you have, type:

```bash
$ xtp2code -v
```

xtp2code converts an XTP Extension Plugin to Go or MoonBit source code
for use with XTP's APIs. It can generate simple custom datatypes and/or Host SDK code
and/or Plugin PDK code. For input, it can process either a schema.yaml file
or it can query the XTP API directly for a given app ID (for the authenticated
user) and process all extension plugin definitions.

Note that if `-appid` is provided, the "XTP_TOKEN" environment variable must
be set for the logged-in XTP user (`xtp auto login`).

Usage:

```
xtp2code \
 -lang=[go|mbt] \
 -pkg=<packageName> \
 [-q ] \
 [-appid=<id> | -yaml=<filename>] \
 [-force] \
 [-host=<filename>] \
 [-plugin=<filename>] \
 [-types=<filename>]
```

[Go]: https://go.dev

## Build Examples

To build all of the examples, `cd` to the `examples` directory and run:

```bash
$ ./build-all.sh
```

To build one of the examples, `cd` to that example subdirectory and run:

```bash
$ ./build.sh
```

## Build from XTP API

To generate code for XTP Extension Plugins directly from the XTP API,
first find the app ID by typing `xtp app list`.

Then run the following to generate Go code:

```bash
$ xtp2code -lang=go -host=api-go-host -plugin=api-go-plugin -types=api-go-types -appid=app_01j1b1mek5frq9x7ymk52m7bw5
```

or run the following to generate MoonBit code:

```bash
$ xtp2code -lang=mbt -plugin=api-mbt-plugin -types=api-mbt-types -appid=app_01j1b1mek5frq9x7ymk52m7bw5
```

## Push and Bind Plugin

Once a plugin has been built successfully, it needs to be pushed to XTP
and then installed (also known as "binding").

To push a plugin to XTP, first build it (see above) and make sure that its
`xtp.toml` file has been filled in correctly with fields:

* `app_id`
* `extension_point_id`
* `name`

You can then run:

```bash
$ xtp plugin push
$ xtp plugin bind
```
