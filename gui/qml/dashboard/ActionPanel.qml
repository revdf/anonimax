import QtQuick
import QtQuick.Layouts

Item {
    id: root

    required property var backend
    property bool safeKill: false
    signal panicRequested()

    readonly property bool iconOnly: width < 700

    implicitHeight: content.implicitHeight

    ColumnLayout {
        id: content

        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        spacing: Theme.spacingSmall

        GridLayout {
            Layout.fillWidth: true
            columns: 5
            columnSpacing: Theme.spacingSmall
            rowSpacing: Theme.spacingSmall

            ActionButton {
                Layout.fillWidth: true
                text: "New Identity"
                iconName: "identity"
                iconOnly: root.iconOnly
                tooltip: "Requests new Tor circuits for new connections. The exit IP may stay the same; existing connections are unchanged."
                enabled: root.backend.statusAvailable && root.backend.anonsurfActive
                         && root.backend.torActive && !root.backend.busy
                onClicked: root.backend.changeIdentity()
            }

            ActionButton {
                Layout.fillWidth: true
                text: "Verify IP"
                iconName: "globe"
                iconOnly: root.iconOnly
                tooltip: "Checks the public IP used by this request and whether it belongs to a known Tor exit node."
                enabled: root.backend.statusAvailable && !root.backend.busy
                onClicked: root.backend.checkIP()
            }

            ActionButton {
                Layout.fillWidth: true
                text: "Restart"
                iconName: "restart"
                iconOnly: root.iconOnly
                tooltip: "Restarts Tor and rebuilds routing from the current Anonimax configuration."
                enabled: root.backend.statusAvailable && root.backend.anonsurfActive
                         && !root.backend.busy
                onClicked: root.backend.restart(root.safeKill)
            }

            ActionButton {
                Layout.fillWidth: true
                text: "Tor Monitor"
                iconName: "monitor"
                iconOnly: root.iconOnly
                tooltip: "Opens nyx, the Tor connection monitor, in a Konsole window."
                enabled: root.backend.statusAvailable && root.backend.anonsurfActive
                         && root.backend.torActive && !root.backend.busy
                onClicked: root.backend.openMonitor()
            }

            ActionButton {
                Layout.fillWidth: true
                text: "Refresh"
                iconName: "refresh"
                iconOnly: root.iconOnly
                tooltip: "Reads the Anonimax status right away instead of waiting for the periodic update."
                enabled: !root.backend.busy
                onClicked: root.backend.refresh()
            }
        }

        GridLayout {
            Layout.fillWidth: true
            columns: 2
            columnSpacing: Theme.spacingSmall
            rowSpacing: Theme.spacingSmall

            ActionButton {
                Layout.fillWidth: true
                text: root.backend.autochangeActive ? "Stop Auto-change" : "Auto-change"
                iconName: "timer"
                iconOnly: root.iconOnly
                tooltip: root.backend.autochangeActive
                         ? "Stop automatic Tor identity rotation."
                         : "Rotate Tor identity automatically every 10 minutes."
                enabled: root.backend.statusAvailable && root.backend.anonsurfActive
                         && root.backend.torActive && !root.backend.busy
                onClicked: root.backend.toggleAutochange()
            }

            ActionButton {
                Layout.fillWidth: true
                text: "Panic"
                iconName: "zap"
                iconOnly: root.iconOnly
                danger: true
                tooltip: "Emergency shutdown: cuts network, stops Tor, clears traces and powers off the machine."
                enabled: !root.backend.busy
                onClicked: root.panicRequested()
            }
        }
    }
}
