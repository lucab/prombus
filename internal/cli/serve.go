package cli

import (
	"bytes"
	"fmt"
	"os"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdServe = &cobra.Command{
		Use:  "serve",
		RunE: runServe,
	}
)

const intro = `
<node>
        <interface name="com.github.lucab.Prombus.Observable">
                <method name="PromMetrics">
                        <arg direction="out" type="a(y)"/>
                </method>
        </interface>` + introspect.IntrospectDataString + `</node> `

type registry bool

func (r registry) PromMetrics() ([]byte, *dbus.Error) {
	logrus.Debug("gathering metrics")
	buf := bytes.Buffer{}

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return nil, dbus.NewError("PromError", []interface{}{err})
	}

	for _, mf := range mfs {
		if _, err := expfmt.MetricFamilyToText(&buf, mf); err != nil {
			return nil, dbus.NewError("PromError", []interface{}{err})
		}
	}

	return buf.Bytes(), nil
}

func runServe(cmd *cobra.Command, cmdArgs []string) error {
	logrus.Debug("setting up DBus endpoint")

	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	reply, err := conn.RequestName("com.github.lucab.Prombus",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	reg := registry(true)
	conn.Export(reg, "/com/github/lucab/Prombus", "com.github.lucab.Prombus.Observable")
	conn.Export(introspect.Introspectable(intro), "/com/github/lucab/Prombus",
		"org.freedesktop.DBus.Introspectable")
	logrus.Debug("listening on com.github.lucab.Prombus")
	select {}

	return nil
}
