# prombus

`prombus` exposes its own metrics, in Prometheus textual format, over a DBus endpoint.

It can be bridged via [`local_exporter`](https://github.com/lucab/local_exporter), with the following configuration:

```toml
[bridge.selectors.prombus]
kind = "dbus"
bus = "system"
destination = "com.github.lucab.Prombus"
method = "com.github.lucab.Prombus.Observable.PromMetrics"
path = "/com/github/lucab/Prombus"
```
