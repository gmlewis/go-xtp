# go-xtp - An experimental Go and MoonBit code generator for XTP Extension Plugins

For more information about XTP, please see:
https://docs.xtp.dylibso.com/docs/overview

## Build

To build all of the examples, `cd` to the `examples` directory and run:

```bash
$ ./build-all.sh
```

To build one of the examples, `cd` to that example subdirectory and run:

```bash
$ ./build.sh
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
