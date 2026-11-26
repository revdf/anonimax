import QtQuick
import QtQuick.Layouts

Rectangle {
    id: root

    required property var config

    implicitHeight: content.implicitHeight + Theme.spacingLarge * 2
    color: Theme.bgPanel
    border.width: 1
    border.color: Theme.borderLight

    ColumnLayout {
        id: content

        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: Theme.spacingLarge
        spacing: Theme.spacingMedium

        Text {
            text: "Privacy options"
            color: Theme.textSecondary
            font.family: Theme.uiFont
            font.pixelSize: Theme.detailSize
            font.letterSpacing: 1
        }

        Rectangle {
            Layout.fillWidth: true
            Layout.preferredHeight: 1
            color: Theme.borderLight
        }

        GridLayout {
            Layout.fillWidth: true
            columns: width < 500 ? 1 : 2
            columnSpacing: Theme.spacingLarge
            rowSpacing: Theme.spacingSmall

            ToggleOption {
                Layout.fillWidth: true
                label: "Randomize MAC"
                tooltip: "Changes the network interface MAC address to a random value at start."
                checked: root.config.changeMac
                enabled: !root.config.busy
                onClicked: root.config.setSetting("change_mac", root.config.changeMac ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Renew local IP"
                tooltip: "Runs dhclient to request a new IP address from the router at start."
                checked: root.config.renewIP
                enabled: !root.config.busy
                onClicked: root.config.setSetting("renew_ip", root.config.renewIP ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Force UTC timezone"
                tooltip: "Sets the system timezone to UTC while Anonimax is active, restoring it on stop."
                checked: root.config.forceUTC
                enabled: !root.config.busy
                onClicked: root.config.setSetting("force_utc", root.config.forceUTC ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Block WebRTC"
                tooltip: "Adds iptables rules to block STUN ports (3478/5349) and prevent WebRTC IP leaks."
                checked: root.config.blockWebRTC
                enabled: !root.config.busy
                onClicked: root.config.setSetting("block_webrtc", root.config.blockWebRTC ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Block telemetry"
                tooltip: "Drops traffic to known Google, Microsoft and Amazon telemetry IP ranges."
                accentColor: Theme.warning
                checked: root.config.blockTelemetry
                enabled: !root.config.busy
                onClicked: root.config.setSetting("block_telemetry", root.config.blockTelemetry ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Disable cam & mic"
                tooltip: "Unloads camera and audio kernel modules at start and reloads them on stop."
                accentColor: Theme.warning
                checked: root.config.disableCameraMic
                enabled: !root.config.busy
                onClicked: root.config.setSetting("disable_camera_mic", root.config.disableCameraMic ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Clean on stop"
                tooltip: "Clears journal, /tmp, DNS cache and browser data via BleachBit when stopping."
                checked: root.config.cleanOnStop
                enabled: !root.config.busy
                onClicked: root.config.setSetting("clean_on_stop", root.config.cleanOnStop ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Use VPN"
                tooltip: "Connects to the VPN defined in vpn_config_path before starting Tor."
                checked: root.config.useVPN
                enabled: !root.config.busy
                onClicked: root.config.setSetting("use_vpn", root.config.useVPN ? "false" : "true")
            }

            ToggleOption {
                Layout.fillWidth: true
                label: "Use proxy"
                tooltip: "Routes traffic through the configured SOCKS5/HTTP proxy."
                checked: root.config.useProxy
                enabled: !root.config.busy
                onClicked: root.config.setSetting("use_proxy", root.config.useProxy ? "false" : "true")
            }
        }
    }
}
